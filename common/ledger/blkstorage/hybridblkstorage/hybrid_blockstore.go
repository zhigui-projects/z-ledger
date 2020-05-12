/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"github.com/hyperledger/fabric-protos-go/ledger/archive"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/ledger"
	"github.com/hyperledger/fabric/common/ledger/blkstorage"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
)

// hybridBlockStore - hybrid implementation for `BlockStore`
type hybridBlockStore struct {
	id      string
	conf    *Conf
	fileMgr *hybridBlockfileMgr
	stats   *ledgerStats
}

// NewHybridBlockStore constructs a `HybridBlockStore`
func newHybridBlockStore(id string, conf *Conf, indexConfig *blkstorage.IndexConfig,
	dbHandle *leveldbhelper.DBHandle, stats *stats) *hybridBlockStore {
	fileMgr := newBlockfileMgr(id, conf, indexConfig, dbHandle)

	// create ledgerStats and initialize blockchain_height stat
	ledgerStats := stats.ledgerStats(id)
	info := fileMgr.getBlockchainInfo()
	ledgerStats.updateBlockchainHeight(info.Height)

	return &hybridBlockStore{id, conf, fileMgr, ledgerStats}
}

func (store *hybridBlockStore) TransferBlockFiles() error {
	return store.fileMgr.transferBlockFiles()
}

func (store *hybridBlockStore) GetArchiveMetaInfo() (*archive.ArchiveMetaInfo, error) {
	return store.fileMgr.loadArchiveMetaInfo()
}

func (store *hybridBlockStore) UpdateArchiveMetaInfo(metaInfo *archive.ArchiveMetaInfo) {
	if err := store.fileMgr.saveArchiveMetaInfo(metaInfo); err != nil {
		logger.Errorf("Saving archive meta info failed with error: %s", err)
		return
	}
	store.fileMgr.updateArchiveMetaInfo(metaInfo)
}

// AddBlock adds a new block
func (store *hybridBlockStore) AddBlock(block *common.Block) error {
	// track elapsed time to collect block commit time
	startBlockCommit := time.Now()
	result := store.fileMgr.addBlock(block)
	elapsedBlockCommit := time.Since(startBlockCommit)

	store.updateBlockStats(block.Header.Number, elapsedBlockCommit)

	return result
}

// GetBlockchainInfo returns the current info about blockchain
func (store *hybridBlockStore) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	return store.fileMgr.getBlockchainInfo(), nil
}

// RetrieveBlocks returns an iterator that can be used for iterating over a range of blocks
func (store *hybridBlockStore) RetrieveBlocks(startNum uint64) (ledger.ResultsIterator, error) {
	return store.fileMgr.retrieveBlocks(startNum)
}

// RetrieveBlockByHash returns the block for given block-hash
func (store *hybridBlockStore) RetrieveBlockByHash(blockHash []byte) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByHash(blockHash)
}

// RetrieveBlockByNumber returns the block at a given blockchain height
func (store *hybridBlockStore) RetrieveBlockByNumber(blockNum uint64) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByNumber(blockNum)
}

// RetrieveTxByID returns a transaction for given transaction id
func (store *hybridBlockStore) RetrieveTxByID(txID string) (*common.Envelope, error) {
	return store.fileMgr.retrieveTransactionByID(txID)
}

// RetrieveTxByID returns a transaction for given transaction id
func (store *hybridBlockStore) RetrieveTxByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error) {
	return store.fileMgr.retrieveTransactionByBlockNumTranNum(blockNum, tranNum)
}

func (store *hybridBlockStore) RetrieveBlockByTxID(txID string) (*common.Block, error) {
	return store.fileMgr.retrieveBlockByTxID(txID)
}

func (store *hybridBlockStore) RetrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	return store.fileMgr.retrieveTxValidationCodeByTxID(txID)
}

// Shutdown shuts down the block store
func (store *hybridBlockStore) Shutdown() {
	logger.Debugf("closing fs blockStore:%s", store.id)
	store.fileMgr.close()
}

func (store *hybridBlockStore) updateBlockStats(blockNum uint64, blockstorageCommitTime time.Duration) {
	store.stats.updateBlockchainHeight(blockNum + 1)
	store.stats.updateBlockstorageCommitTime(blockstorageCommitTime)
}
