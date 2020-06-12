/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/ledger/archive"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger/blkstorage"
	"github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/core/ledger/archive"
	"github.com/hyperledger/fabric/core/ledger/archive/eventbus"
	"github.com/hyperledger/fabric/core/ledger/dfs"
	dc "github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var logger = flogging.MustGetLogger("hybridblkstorage")

const (
	blockfilePrefix = "blockfile_"
)

var (
	blkMgrInfoKey      = []byte("blkMgrInfo")
	archiveMetaInfoKey = []byte("archiveMetaInfo")
)

type hybridBlockfileMgr struct {
	rootDir           string
	conf              *Conf
	db                *leveldbhelper.DBHandle
	index             index
	cpInfo            *checkpointInfo
	cpInfoCond        *sync.Cond
	currentFileWriter *blockfileWriter
	bcInfo            atomic.Value
	amInfo            *pb.ArchiveMetaInfo
	amInfoCond        *sync.Cond
	dfsClient         dc.FsClient
	lock              sync.RWMutex
}

/*
Creates a new manager that will manage the files used for block persistence.
This manager manages the file system FS including
  -- the directory where the files are stored
  -- the individual files where the blocks are stored
  -- the checkpoint which tracks the latest file being persisted to
  -- the index which tracks what block and transaction is in what file
When a new blockfile manager is started (i.e. only on start-up), it checks
if this start-up is the first time the system is coming up or is this a restart
of the system.

The blockfile manager stores blocks of data into a file system.  That file
storage is done by creating sequentially numbered files of a configured size
i.e blockfile_000000, blockfile_000001, etc..

Each transaction in a block is stored with information about the number of
bytes in that transaction
 Adding txLoc [fileSuffixNum=0, offset=3, bytesLength=104] for tx [1:0] to index
 Adding txLoc [fileSuffixNum=0, offset=107, bytesLength=104] for tx [1:1] to index
Each block is stored with the total encoded length of that block as well as the
tx location offsets.

Remember that these steps are only done once at start-up of the system.
At start up a new manager:
  *) Checks if the directory for storing files exists, if not creates the dir
  *) Checks if the key value database exists, if not creates one
       (will create a db dir)
  *) Determines the checkpoint information (cpinfo) used for storage
		-- Loads from db if exist, if not instantiate a new cpinfo
		-- If cpinfo was loaded from db, compares to FS
		-- If cpinfo and file system are not in sync, syncs cpInfo from FS
  *) Starts a new file writer
		-- truncates file per cpinfo to remove any excess past last block
  *) Determines the index information used to find tx and blocks in
  the file blkstorage
		-- Instantiates a new blockIdxInfo
		-- Loads the index from the db if exists
		-- syncIndex comparing the last block indexed to what is in the FS
		-- If index and file system are not in sync, syncs index from the FS
  *)  Updates blockchain info used by the APIs
*/
func newBlockfileMgr(id string, conf *Conf, indexConfig *blkstorage.IndexConfig,
	indexStore *leveldbhelper.DBHandle, archiveConf *archive.Config) *hybridBlockfileMgr {
	logger.Debugf("newBlockfileMgr() initializing hybrid-file-based block storage for ledger: %s ", id)
	//Determine the root directory for the blockfile storage, if it does not exist create it
	rootDir := conf.getLedgerBlockDir(id)
	_, err := util.CreateDirIfMissing(rootDir)
	if err != nil {
		panic(fmt.Sprintf("Error creating block storage root dir [%s]: %s", rootDir, err))
	}
	dfsConf := &dc.Config{Type: archiveConf.Type, HdfsConf: archiveConf.HdfsConf, IpfsConf: archiveConf.IpfsConf}
	client, err := dfs.NewDfsClient(dfsConf)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not connect to HDFS, due to %+v", err))
	}
	// Instantiate the manager, i.e. blockFileMgr structure
	mgr := &hybridBlockfileMgr{rootDir: rootDir, conf: conf, db: indexStore, dfsClient: client}

	amInfo, err := mgr.loadArchiveMetaInfo()
	if err != nil {
		logger.Errorf("Could not get archive meta info for ledger: %s from db: %+v", id, err)
	}
	if amInfo == nil {
		logger.Info("Getting archive info from dfs storage")
		if amInfo, err = constructArchiveMetaInfoFromDfsBlockFiles(rootDir, client); err != nil {
			logger.Errorf("Could not build archive meta info from dfs block files: %+v", err)
		}
		logger.Infof("Archive meta info constructed by scanning the dfs blocks dir: %s", spew.Sdump(amInfo))
	} else {
		logger.Info("Syncing archive meta info from dfs (if needed)")
		syncArchiveMetaInfoFromDfs(rootDir, amInfo, client)
	}
	amInfo.ChannelId = id
	mgr.amInfo = amInfo
	mgr.amInfoCond = sync.NewCond(&sync.Mutex{})
	if err := mgr.saveArchiveMetaInfo(amInfo); err != nil {
		logger.Errorf("Could not save archive meta info to db: %s", err)
	}

	if err := eventbus.Get(id).Subscribe(eventbus.ArchiveByTxDate, mgr.archiveFn); err != nil {
		logger.Errorf("Could not subscribe to archive event for chain[id=%s]", id)
	}

	// cp = checkpointInfo, retrieve from the database the file suffix or number of where blocks were stored.
	// It also retrieves the current size of that file and the last block number that was written to that file.
	// At init checkpointInfo:latestFileChunkSuffixNum=[0], latestFileChunksize=[0], lastBlockNumber=[0]
	cpInfo, err := mgr.loadCurrentInfo()
	if err != nil {
		panic(fmt.Sprintf("Could not get block file info for current block file from db: %s", err))
	}
	if cpInfo == nil {
		logger.Info(`Getting block information from block storage`)
		if cpInfo, err = constructCheckpointInfoFromBlockFiles(rootDir); err != nil {
			panic(fmt.Sprintf("Could not build checkpoint info from block files: %s", err))
		}
		logger.Debugf("Info constructed by scanning the blocks dir = %s", spew.Sdump(cpInfo))
	} else {
		logger.Debug(`Synching block information from block storage (if needed)`)
		syncCPInfoFromFS(rootDir, cpInfo)
	}
	err = mgr.saveCurrentInfo(cpInfo, true)
	if err != nil {
		panic(fmt.Sprintf("Could not save next block file info to db: %s", err))
	}

	//Open a writer to the file identified by the number and truncate it to only contain the latest block
	// that was completely saved (file system, index, cpinfo, etc)
	currentFileWriter, err := newBlockfileWriter(deriveBlockfilePath(rootDir, cpInfo.latestFileChunkSuffixNum))
	if err != nil {
		panic(fmt.Sprintf("Could not open writer to current file: %s", err))
	}
	//Truncate the file to remove excess past last block
	err = currentFileWriter.truncateFile(cpInfo.latestFileChunksize)
	if err != nil {
		panic(fmt.Sprintf("Could not truncate current file to known size in db: %s", err))
	}

	// Create a new KeyValue store database handler for the blocks index in the key value database
	if mgr.index, err = newBlockIndex(indexConfig, indexStore); err != nil {
		panic(fmt.Sprintf("error in block index: %s", err))
	}

	// Update the manager with the checkpoint info and the file writer
	mgr.cpInfo = cpInfo
	mgr.currentFileWriter = currentFileWriter
	// Create a checkpoint condition (event) variable, for the  goroutine waiting for
	// or announcing the occurrence of an event.
	mgr.cpInfoCond = sync.NewCond(&sync.Mutex{})

	// init BlockchainInfo for external API's
	bcInfo := &common.BlockchainInfo{
		Height:            0,
		CurrentBlockHash:  nil,
		PreviousBlockHash: nil}

	if !cpInfo.isChainEmpty {
		//If start up is a restart of an existing storage, sync the index from block storage and update BlockchainInfo for external API's
		mgr.syncIndex()
		lastBlockHeader, err := mgr.retrieveBlockHeaderByNumber(cpInfo.lastBlockNumber)
		if err != nil {
			panic(fmt.Sprintf("Could not retrieve header of the last block form file: %s", err))
		}
		lastBlockHash := protoutil.BlockHeaderHash(lastBlockHeader)
		previousBlockHash := lastBlockHeader.PreviousHash
		bcInfo = &common.BlockchainInfo{
			Height:            cpInfo.lastBlockNumber + 1,
			CurrentBlockHash:  lastBlockHash,
			PreviousBlockHash: previousBlockHash}
	}
	mgr.bcInfo.Store(bcInfo)
	return mgr
}

//cp = checkpointInfo, from the database gets the file suffix and the size of
// the file of where the last block was written.  Also retrieves contains the
// last block number that was written.  At init
//checkpointInfo:latestFileChunkSuffixNum=[0], latestFileChunksize=[0], lastBlockNumber=[0]
func syncCPInfoFromFS(rootDir string, cpInfo *checkpointInfo) {
	logger.Debugf("Starting checkpoint=%s", cpInfo)
	//Checks if the file suffix of where the last block was written exists
	filePath := deriveBlockfilePath(rootDir, cpInfo.latestFileChunkSuffixNum)
	exists, size, err := util.FileExists(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error in checking whether file [%s] exists: %s", filePath, err))
	}
	logger.Debugf("status of file [%s]: exists=[%t], size=[%d]", filePath, exists, size)
	//Test is !exists because when file number is first used the file does not exist yet
	//checks that the file exists and that the size of the file is what is stored in cpinfo
	//status of file [/tmp/tests/ledger/blkstorage/fsblkstorage/blocks/blockfile_000000]: exists=[false], size=[0]
	if !exists || int(size) == cpInfo.latestFileChunksize {
		// check point info is in sync with the file on disk
		return
	}
	//Scan the file system to verify that the checkpoint info stored in db is correct
	_, endOffsetLastBlock, numBlocks, err := scanForLastCompleteBlock(
		rootDir, cpInfo.latestFileChunkSuffixNum, int64(cpInfo.latestFileChunksize))
	if err != nil {
		panic(fmt.Sprintf("Could not open current file for detecting last block in the file: %s", err))
	}
	cpInfo.latestFileChunksize = int(endOffsetLastBlock)
	if numBlocks == 0 {
		return
	}
	//Updates the checkpoint info for the actual last block number stored and it's end location
	if cpInfo.isChainEmpty {
		cpInfo.lastBlockNumber = uint64(numBlocks - 1)
	} else {
		cpInfo.lastBlockNumber += uint64(numBlocks)
	}
	cpInfo.isChainEmpty = false
	logger.Debugf("Checkpoint after updates by scanning the last file segment:%s", cpInfo)
}

func deriveBlockfilePath(rootDir string, suffixNum int) string {
	return rootDir + "/" + blockfilePrefix + fmt.Sprintf("%06d", suffixNum)
}

func (mgr *hybridBlockfileMgr) close() {
	mgr.currentFileWriter.close()
}

func (mgr *hybridBlockfileMgr) moveToNextFile() {
	cpInfo := &checkpointInfo{
		latestFileChunkSuffixNum: mgr.cpInfo.latestFileChunkSuffixNum + 1,
		latestFileChunksize:      0,
		lastBlockNumber:          mgr.cpInfo.lastBlockNumber}

	nextFileWriter, err := newBlockfileWriter(
		deriveBlockfilePath(mgr.rootDir, cpInfo.latestFileChunkSuffixNum))

	if err != nil {
		panic(fmt.Sprintf("Could not open writer to next file: %s", err))
	}
	mgr.currentFileWriter.close()
	err = mgr.saveCurrentInfo(cpInfo, true)
	if err != nil {
		panic(fmt.Sprintf("Could not save next block file info to db: %s", err))
	}
	mgr.currentFileWriter = nextFileWriter
	mgr.updateCheckpoint(cpInfo)
}

func (mgr *hybridBlockfileMgr) addBlock(block *common.Block) error {
	bcInfo := mgr.getBlockchainInfo()
	if block.Header.Number != bcInfo.Height {
		return errors.Errorf(
			"block number should have been %d but was %d",
			mgr.getBlockchainInfo().Height, block.Header.Number,
		)
	}

	// Add the previous hash check - Though, not essential but may not be a bad idea to
	// verify the field `block.Header.PreviousHash` present in the block.
	// This check is a simple bytes comparison and hence does not cause any observable performance penalty
	// and may help in detecting a rare scenario if there is any bug in the ordering service.
	if !bytes.Equal(block.Header.PreviousHash, bcInfo.CurrentBlockHash) {
		return errors.Errorf(
			"unexpected Previous block hash. Expected PreviousHash = [%x], PreviousHash referred in the latest block= [%x]",
			bcInfo.CurrentBlockHash, block.Header.PreviousHash,
		)
	}
	blockBytes, info, err := serializeBlock(block)
	if err != nil {
		return errors.WithMessage(err, "error serializing block")
	}
	blockHash := protoutil.BlockHeaderHash(block.Header)
	//Get the location / offset where each transaction starts in the block and where the block ends
	txOffsets := info.txOffsets
	currentOffset := mgr.cpInfo.latestFileChunksize

	blockBytesLen := len(blockBytes)
	blockBytesEncodedLen := proto.EncodeVarint(uint64(blockBytesLen))
	totalBytesToAppend := blockBytesLen + len(blockBytesEncodedLen)

	//Determine if we need to start a new file since the size of this block
	//exceeds the amount of space left in the current file
	if currentOffset+totalBytesToAppend > mgr.conf.maxBlockfileSize {
		mgr.moveToNextFile()
		currentOffset = 0
	}
	//append blockBytesEncodedLen to the file
	err = mgr.currentFileWriter.append(blockBytesEncodedLen, false)
	if err == nil {
		//append the actual block bytes to the file
		err = mgr.currentFileWriter.append(blockBytes, true)
	}
	if err != nil {
		truncateErr := mgr.currentFileWriter.truncateFile(mgr.cpInfo.latestFileChunksize)
		if truncateErr != nil {
			panic(fmt.Sprintf("Could not truncate current file to known size after an error during block append: %s", err))
		}
		return errors.WithMessage(err, "error appending block to file")
	}

	//Update the checkpoint info with the results of adding the new block
	currentCPInfo := mgr.cpInfo
	newCPInfo := &checkpointInfo{
		latestFileChunkSuffixNum: currentCPInfo.latestFileChunkSuffixNum,
		latestFileChunksize:      currentCPInfo.latestFileChunksize + totalBytesToAppend,
		isChainEmpty:             false,
		lastBlockNumber:          block.Header.Number}
	//save the checkpoint information in the database
	if err = mgr.saveCurrentInfo(newCPInfo, false); err != nil {
		truncateErr := mgr.currentFileWriter.truncateFile(currentCPInfo.latestFileChunksize)
		if truncateErr != nil {
			panic(fmt.Sprintf("Error in truncating current file to known size after an error in saving checkpoint info: %s", err))
		}
		return errors.WithMessage(err, "error saving current file info to db")
	}

	//Index block file location pointer updated with file suffex and offset for the new block
	blockFLP := &fileLocPointer{fileSuffixNum: newCPInfo.latestFileChunkSuffixNum}
	blockFLP.offset = currentOffset
	// shift the txoffset because we prepend length of bytes before block bytes
	for _, txOffset := range txOffsets {
		txOffset.loc.offset += len(blockBytesEncodedLen)
	}
	//save the index in the database
	if err = mgr.index.indexBlock(&blockIdxInfo{
		blockNum: block.Header.Number, blockHash: blockHash,
		flp: blockFLP, txOffsets: txOffsets, metadata: block.Metadata}); err != nil {
		return err
	}

	//update the checkpoint info (for storage) and the blockchain info (for APIs) in the manager
	mgr.updateCheckpoint(newCPInfo)
	mgr.updateBlockchainInfo(blockHash, block)
	return nil
}

func (mgr *hybridBlockfileMgr) syncIndex() error {
	var lastBlockIndexed uint64
	var indexEmpty bool
	var err error
	//from the database, get the last block that was indexed
	if lastBlockIndexed, err = mgr.index.getLastBlockIndexed(); err != nil {
		if err != errIndexEmpty {
			return err
		}
		indexEmpty = true
	}

	//initialize index to file number:zero, offset:zero and blockNum:0
	startFileNum := 0
	startOffset := 0
	skipFirstBlock := false
	//get the last file that blocks were added to using the checkpoint info
	endFileNum := mgr.cpInfo.latestFileChunkSuffixNum
	startingBlockNum := uint64(0)

	//if the index stored in the db has value, update the index information with those values
	if !indexEmpty {
		if lastBlockIndexed == mgr.cpInfo.lastBlockNumber {
			logger.Debug("Both the block files and indices are in sync.")
			return nil
		}
		logger.Debugf("Last block indexed [%d], Last block present in block files [%d]", lastBlockIndexed, mgr.cpInfo.lastBlockNumber)
		var flp *fileLocPointer
		if flp, err = mgr.index.getBlockLocByBlockNum(lastBlockIndexed); err != nil {
			return err
		}
		startFileNum = flp.fileSuffixNum
		startOffset = flp.locPointer.offset
		skipFirstBlock = true
		startingBlockNum = lastBlockIndexed + 1
	} else {
		logger.Debugf("No block indexed, Last block present in block files=[%d]", mgr.cpInfo.lastBlockNumber)
	}

	logger.Infof("Start building index from block [%d] to last block [%d]", startingBlockNum, mgr.cpInfo.lastBlockNumber)

	//open a blockstream to the file location that was stored in the index
	var stream *blockStream
	if stream, err = newBlockStream(mgr.rootDir, startFileNum, int64(startOffset), endFileNum); err != nil {
		return err
	}
	var blockBytes []byte
	var blockPlacementInfo *blockPlacementInfo

	if skipFirstBlock {
		if blockBytes, _, err = stream.nextBlockBytesAndPlacementInfo(); err != nil {
			return err
		}
		if blockBytes == nil {
			return errors.Errorf("block bytes for block num = [%d] should not be nil here. The indexes for the block are already present",
				lastBlockIndexed)
		}
	}

	//Should be at the last block already, but go ahead and loop looking for next blockBytes.
	//If there is another block, add it to the index.
	//This will ensure block indexes are correct, for example if peer had crashed before indexes got updated.
	blockIdxInfo := &blockIdxInfo{}
	for {
		if blockBytes, blockPlacementInfo, err = stream.nextBlockBytesAndPlacementInfo(); err != nil {
			return err
		}
		if blockBytes == nil {
			break
		}
		info, err := extractSerializedBlockInfo(blockBytes)
		if err != nil {
			return err
		}

		//The blockStartOffset will get applied to the txOffsets prior to indexing within indexBlock(),
		//therefore just shift by the difference between blockBytesOffset and blockStartOffset
		numBytesToShift := int(blockPlacementInfo.blockBytesOffset - blockPlacementInfo.blockStartOffset)
		for _, offset := range info.txOffsets {
			offset.loc.offset += numBytesToShift
		}

		//Update the blockIndexInfo with what was actually stored in file system
		blockIdxInfo.blockHash = protoutil.BlockHeaderHash(info.blockHeader)
		blockIdxInfo.blockNum = info.blockHeader.Number
		blockIdxInfo.flp = &fileLocPointer{fileSuffixNum: blockPlacementInfo.fileNum,
			locPointer: locPointer{offset: int(blockPlacementInfo.blockStartOffset)}}
		blockIdxInfo.txOffsets = info.txOffsets
		blockIdxInfo.metadata = info.metadata

		logger.Debugf("syncIndex() indexing block [%d]", blockIdxInfo.blockNum)
		if err = mgr.index.indexBlock(blockIdxInfo); err != nil {
			return err
		}
		if blockIdxInfo.blockNum%10000 == 0 {
			logger.Infof("Indexed block number [%d]", blockIdxInfo.blockNum)
		}
	}
	logger.Infof("Finished building index. Last block indexed [%d]", blockIdxInfo.blockNum)
	return nil
}

func (mgr *hybridBlockfileMgr) getBlockchainInfo() *common.BlockchainInfo {
	return mgr.bcInfo.Load().(*common.BlockchainInfo)
}

func (mgr *hybridBlockfileMgr) updateCheckpoint(cpInfo *checkpointInfo) {
	mgr.cpInfoCond.L.Lock()
	defer mgr.cpInfoCond.L.Unlock()
	mgr.cpInfo = cpInfo
	logger.Debugf("Broadcasting about update checkpointInfo: %s", cpInfo)
	mgr.cpInfoCond.Broadcast()
}

func (mgr *hybridBlockfileMgr) updateBlockchainInfo(latestBlockHash []byte, latestBlock *common.Block) {
	currentBCInfo := mgr.getBlockchainInfo()
	newBCInfo := &common.BlockchainInfo{
		Height:            currentBCInfo.Height + 1,
		CurrentBlockHash:  latestBlockHash,
		PreviousBlockHash: latestBlock.Header.PreviousHash}

	mgr.bcInfo.Store(newBCInfo)
}

func (mgr *hybridBlockfileMgr) retrieveBlockByHash(blockHash []byte) (*common.Block, error) {
	logger.Debugf("retrieveBlockByHash() - blockHash = [%#v]", blockHash)
	loc, err := mgr.index.getBlockLocByHash(blockHash)
	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *hybridBlockfileMgr) retrieveBlockByNumber(blockNum uint64) (*common.Block, error) {
	logger.Debugf("retrieveBlockByNumber() - blockNum = [%d]", blockNum)

	// interpret math.MaxUint64 as a request for last block
	if blockNum == math.MaxUint64 {
		blockNum = mgr.getBlockchainInfo().Height - 1
	}

	loc, err := mgr.index.getBlockLocByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *hybridBlockfileMgr) retrieveBlockByTxID(txID string) (*common.Block, error) {
	logger.Debugf("retrieveBlockByTxID() - txID = [%s]", txID)

	loc, err := mgr.index.getBlockLocByTxID(txID)

	if err != nil {
		return nil, err
	}
	return mgr.fetchBlock(loc)
}

func (mgr *hybridBlockfileMgr) retrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	logger.Debugf("retrieveTxValidationCodeByTxID() - txID = [%s]", txID)
	return mgr.index.getTxValidationCodeByTxID(txID)
}

func (mgr *hybridBlockfileMgr) retrieveBlockHeaderByNumber(blockNum uint64) (*common.BlockHeader, error) {
	logger.Debugf("retrieveBlockHeaderByNumber() - blockNum = [%d]", blockNum)
	loc, err := mgr.index.getBlockLocByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	blockBytes, err := mgr.fetchBlockBytes(loc)
	if err != nil {
		return nil, err
	}
	info, err := extractSerializedBlockInfo(blockBytes)
	if err != nil {
		return nil, err
	}
	return info.blockHeader, nil
}

func (mgr *hybridBlockfileMgr) retrieveBlocks(startNum uint64) (*blocksItr, error) {
	return newBlockItr(mgr, startNum), nil
}

func (mgr *hybridBlockfileMgr) retrieveTransactionByID(txID string) (*common.Envelope, error) {
	logger.Debugf("retrieveTransactionByID() - txId = [%s]", txID)
	loc, err := mgr.index.getTxLoc(txID)
	if err != nil {
		return nil, err
	}
	return mgr.fetchTransactionEnvelope(loc)
}

func (mgr *hybridBlockfileMgr) retrieveTransactionByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error) {
	logger.Debugf("retrieveTransactionByBlockNumTranNum() - blockNum = [%d], tranNum = [%d]", blockNum, tranNum)
	loc, err := mgr.index.getTXLocByBlockNumTranNum(blockNum, tranNum)
	if err != nil {
		return nil, err
	}
	return mgr.fetchTransactionEnvelope(loc)
}

func (mgr *hybridBlockfileMgr) fetchBlock(lp *fileLocPointer) (*common.Block, error) {
	blockBytes, err := mgr.fetchBlockBytes(lp)
	if err != nil {
		return nil, err
	}
	block, err := deserializeBlock(blockBytes)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (mgr *hybridBlockfileMgr) fetchTransactionEnvelope(lp *fileLocPointer) (*common.Envelope, error) {
	logger.Debugf("Entering fetchTransactionEnvelope() %v\n", lp)
	var err error
	var txEnvelopeBytes []byte
	if txEnvelopeBytes, err = mgr.fetchRawBytes(lp); err != nil {
		return nil, err
	}
	_, n := proto.DecodeVarint(txEnvelopeBytes)
	return protoutil.GetEnvelopeFromBlock(txEnvelopeBytes[n:])
}

func (mgr *hybridBlockfileMgr) fetchBlockBytes(lp *fileLocPointer) ([]byte, error) {
	var err error
	var stream hybridBlockStream
	if int32(lp.fileSuffixNum) <= mgr.amInfo.LastArchiveFileSuffix {
		stream, err = newDfsBlockfileStream(mgr.rootDir, lp.fileSuffixNum, int64(lp.offset), mgr.dfsClient, mgr.amInfo.FileProofs[int32(lp.fileSuffixNum)])

	} else {
		stream, err = newBlockfileStream(mgr.rootDir, lp.fileSuffixNum, int64(lp.offset))
	}
	if err != nil {
		return nil, err
	}
	defer stream.close()
	b, err := stream.nextBlockBytes()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (mgr *hybridBlockfileMgr) fetchRawBytes(lp *fileLocPointer) ([]byte, error) {
	var reader fileReader
	var err error
	filePath := deriveBlockfilePath(mgr.rootDir, lp.fileSuffixNum)
	if int32(lp.fileSuffixNum) <= mgr.amInfo.LastArchiveFileSuffix {
		reader, err = newDfsBlockfileReader(filePath, mgr.dfsClient, mgr.amInfo.FileProofs[int32(lp.fileSuffixNum)])
	} else {
		reader, err = newBlockfileReader(filePath)
	}
	if err != nil {
		return nil, err
	}
	defer reader.close()
	b, err := reader.readAt(lp.offset, lp.bytesLength)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Get the current archive meta info that is stored in the database
func (mgr *hybridBlockfileMgr) loadArchiveMetaInfo() (*pb.ArchiveMetaInfo, error) {
	var b []byte
	var err error
	if b, err = mgr.db.Get(archiveMetaInfoKey); b == nil || err != nil {
		return nil, err
	}
	archiveMetaInfo := &pb.ArchiveMetaInfo{}
	if err := proto.Unmarshal(b, archiveMetaInfo); err != nil {
		return nil, err
	}
	logger.Infof("load archiveMetaInfo:%s", spew.Sdump(archiveMetaInfo))
	return archiveMetaInfo, nil
}

func (mgr *hybridBlockfileMgr) saveArchiveMetaInfo(amInfo *pb.ArchiveMetaInfo) error {
	logger.Infof("Saving archive meta info: %s", spew.Sdump(amInfo))
	b, err := proto.Marshal(amInfo)
	if err != nil {
		logger.Errorf("Marshal archive meta info with error: %s", err)
		return err
	}
	if err = mgr.db.Put(archiveMetaInfoKey, b, true); err != nil {
		logger.Errorf("Save archive meta info with error: %s", err)
		return err
	}
	return nil
}

func (mgr *hybridBlockfileMgr) updateArchiveMetaInfo(amInfo *pb.ArchiveMetaInfo) {
	logger.Infof("Updating archive meta info: %s", spew.Sdump(amInfo))
	mgr.amInfoCond.L.Lock()
	defer mgr.amInfoCond.L.Unlock()
	mgr.amInfo = amInfo
	logger.Infof("Broadcasting about update archive meta info: %s", spew.Sdump(amInfo))
	mgr.amInfoCond.Broadcast()
}

func (mgr *hybridBlockfileMgr) transferBlockFiles() error {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	logger.Infof("Transferring block files to dfs if needed")
	lastSentFileNum := int(mgr.amInfo.LastSentFileSuffix)
	latestFileNum := mgr.cpInfo.latestFileChunkSuffixNum
	for ; latestFileNum >= lastSentFileNum+2; lastSentFileNum++ {
		logger.Infof("Start transferring the blockfile[%d] to dfs", lastSentFileNum+1)
		filePath := deriveBlockfilePath(mgr.rootDir, lastSentFileNum+1)
		if _, notExistErr := mgr.dfsClient.Stat(filePath); notExistErr != nil {
			logger.Infof("Blockfile[%s] not exits in dfs, error: %+v", filePath, notExistErr)
			if err := mgr.dfsClient.CopyToRemote(filePath, filePath); err != nil {
				logger.Errorf("Transferring blockfile[%s] failed with error: %+v", filePath, err)
				return err
			}
		} else {
			logger.Warnf("Blockfile already exists[%s] in dfs", filePath)
		}
		if int32(lastSentFileNum+1) > mgr.amInfo.LastSentFileSuffix {
			newAmInfo := &pb.ArchiveMetaInfo{
				LastSentFileSuffix:    int32(lastSentFileNum + 1),
				LastArchiveFileSuffix: mgr.amInfo.LastArchiveFileSuffix,
				FileProofs:            mgr.amInfo.FileProofs,
				ChannelId:             mgr.amInfo.ChannelId,
			}
			if err := mgr.saveArchiveMetaInfo(newAmInfo); err != nil {
				panic(fmt.Sprintf("Could not save archive meta info: %s to db: %+v", spew.Sdump(newAmInfo), err))
			}
			mgr.updateArchiveMetaInfo(newAmInfo)
		}
	}
	return nil
}

func (mgr *hybridBlockfileMgr) archiveFn(channelId string, dateStr string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	logger.Infof("Archiving block files which contains tx before date[%s]", dateStr)
	lastSentFileNum := int(mgr.amInfo.LastSentFileSuffix)
	lastArchiveFileNum := int(mgr.amInfo.LastArchiveFileSuffix)

	var blkLoc *fileLocPointer
	var err error
	if blkLoc, err = mgr.index.getBlockLocByTxDate(dateStr); err != nil {
		logger.Errorf("The block location of given date[%s] has not been found in the index", dateStr)
		return
	}

	// keep the block file contains the tx of given tx date
	fileNum := blkLoc.fileSuffixNum - 1

	if fileNum > lastSentFileNum {
		logger.Errorf("The block files of given date[%s] has not been transferred to dfs yet", dateStr)
		return
	}
	if fileNum <= lastArchiveFileNum {
		logger.Warnf("The block files of given date[%s] has been archived already", dateStr)
		return
	}

	filesInfo, err := ioutil.ReadDir(mgr.conf.getLedgerBlockDir(channelId))
	if err != nil {
		logger.Errorf("archiveFn read dir got error: %s", err)
		return
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Infof("archiveFn skipping File name = %s", name)
			continue
		}
		fileSuffix := strings.TrimPrefix(name, blockfilePrefix)
		suffix, err := strconv.Atoi(fileSuffix)
		if err != nil {
			logger.Errorf("archiveFn convert file suffix[%s] to number got error: %s", fileSuffix, err)
			continue
		}
		if suffix < fileNum {
			// delete the block file
			if err := os.Remove(name); err != nil {
				logger.Errorf("archiveFn remove blockfile[%s] got error: %s", name, err)
				continue
			}
		}
	}

	// update archive meta info
	newAmInfo := &pb.ArchiveMetaInfo{
		LastSentFileSuffix:    mgr.amInfo.LastSentFileSuffix,
		LastArchiveFileSuffix: int32(fileNum),
		FileProofs:            mgr.amInfo.FileProofs,
		ChannelId:             mgr.amInfo.ChannelId,
	}
	if err := mgr.saveArchiveMetaInfo(newAmInfo); err != nil {
		logger.Errorf("Could not save archive meta info: %s to db: %+v", spew.Sdump(newAmInfo), err)
		return
	}
	mgr.updateArchiveMetaInfo(newAmInfo)
}

//Get the current checkpoint information that is stored in the database
func (mgr *hybridBlockfileMgr) loadCurrentInfo() (*checkpointInfo, error) {
	var b []byte
	var err error
	if b, err = mgr.db.Get(blkMgrInfoKey); b == nil || err != nil {
		return nil, err
	}
	i := &checkpointInfo{}
	if err = i.unmarshal(b); err != nil {
		return nil, err
	}
	logger.Debugf("loaded checkpointInfo:%s", i)
	return i, nil
}

func (mgr *hybridBlockfileMgr) saveCurrentInfo(i *checkpointInfo, sync bool) error {
	b, err := i.marshal()
	if err != nil {
		return err
	}
	if err = mgr.db.Put(blkMgrInfoKey, b, sync); err != nil {
		return err
	}
	return nil
}

// scanForLastCompleteBlock scan a given block file and detects the last offset in the file
// after which there may lie a block partially written (towards the end of the file in a crash scenario).
func scanForLastCompleteBlock(rootDir string, fileNum int, startingOffset int64) ([]byte, int64, int, error) {
	//scan the passed file number suffix starting from the passed offset to find the last completed block
	numBlocks := 0
	var lastBlockBytes []byte
	blockStream, errOpen := newBlockfileStream(rootDir, fileNum, startingOffset)
	if errOpen != nil {
		return nil, 0, 0, errOpen
	}
	defer blockStream.close()
	var errRead error
	var blockBytes []byte
	for {
		blockBytes, errRead = blockStream.nextBlockBytes()
		if blockBytes == nil || errRead != nil {
			break
		}
		lastBlockBytes = blockBytes
		numBlocks++
	}
	if errRead == ErrUnexpectedEndOfBlockfile {
		logger.Debugf(`Error:%s
		The error may happen if a crash has happened during block appending.
		Resetting error to nil and returning current offset as a last complete block's end offset`, errRead)
		errRead = nil
	}
	logger.Debugf("scanForLastCompleteBlock(): last complete block ends at offset=[%d]", blockStream.currentOffset)
	return lastBlockBytes, blockStream.currentOffset, numBlocks, errRead
}

// checkpointInfo
type checkpointInfo struct {
	latestFileChunkSuffixNum int
	latestFileChunksize      int
	isChainEmpty             bool
	lastBlockNumber          uint64
}

func (i *checkpointInfo) marshal() ([]byte, error) {
	buffer := proto.NewBuffer([]byte{})
	var err error
	if err = buffer.EncodeVarint(uint64(i.latestFileChunkSuffixNum)); err != nil {
		return nil, errors.Wrapf(err, "error encoding the latestFileChunkSuffixNum [%d]", i.latestFileChunkSuffixNum)
	}
	if err = buffer.EncodeVarint(uint64(i.latestFileChunksize)); err != nil {
		return nil, errors.Wrapf(err, "error encoding the latestFileChunksize [%d]", i.latestFileChunksize)
	}
	if err = buffer.EncodeVarint(i.lastBlockNumber); err != nil {
		return nil, errors.Wrapf(err, "error encoding the lastBlockNumber [%d]", i.lastBlockNumber)
	}
	var chainEmptyMarker uint64
	if i.isChainEmpty {
		chainEmptyMarker = 1
	}
	if err = buffer.EncodeVarint(chainEmptyMarker); err != nil {
		return nil, errors.Wrapf(err, "error encoding chainEmptyMarker [%d]", chainEmptyMarker)
	}
	return buffer.Bytes(), nil
}

func (i *checkpointInfo) unmarshal(b []byte) error {
	buffer := proto.NewBuffer(b)
	var val uint64
	var chainEmptyMarker uint64
	var err error

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.latestFileChunkSuffixNum = int(val)

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.latestFileChunksize = int(val)

	if val, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.lastBlockNumber = val
	if chainEmptyMarker, err = buffer.DecodeVarint(); err != nil {
		return err
	}
	i.isChainEmpty = chainEmptyMarker == 1
	return nil
}

func (i *checkpointInfo) String() string {
	return fmt.Sprintf("latestFileChunkSuffixNum=[%d], latestFileChunksize=[%d], isChainEmpty=[%t], lastBlockNumber=[%d]",
		i.latestFileChunkSuffixNum, i.latestFileChunksize, i.isChainEmpty, i.lastBlockNumber)
}
