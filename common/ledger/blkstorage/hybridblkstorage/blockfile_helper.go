/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/archive"
	la "github.com/hyperledger/fabric/core/ledger/archive"
	dc "github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/pkg/errors"
)

// constructCheckpointInfoFromBlockFiles scans the last blockfile (if any) and construct the checkpoint info
// if the last file contains no block or only a partially written block (potentially because of a crash while writing block to the file),
// this scans the second last file (if any)
func constructCheckpointInfoFromBlockFiles(rootDir string) (*checkpointInfo, error) {
	logger.Debugf("Retrieving checkpoint info from block files")
	var lastFileNum int
	var numBlocksInFile int
	var endOffsetLastBlock int64
	var lastBlockNumber uint64

	var lastBlockBytes []byte
	var lastBlock *common.Block
	var err error

	if lastFileNum, err = retrieveLastFileSuffix(rootDir); err != nil {
		return nil, err
	}
	logger.Debugf("Last file number found = %d", lastFileNum)

	if lastFileNum == -1 {
		cpInfo := &checkpointInfo{0, 0, true, 0}
		logger.Debugf("No block file found")
		return cpInfo, nil
	}

	fileInfo := getFileInfoOrPanic(rootDir, lastFileNum)
	logger.Debugf("Last Block file info: FileName=[%s], FileSize=[%d]", fileInfo.Name(), fileInfo.Size())
	if lastBlockBytes, endOffsetLastBlock, numBlocksInFile, err = scanForLastCompleteBlock(rootDir, lastFileNum, 0); err != nil {
		logger.Errorf("Error scanning last file [num=%d]: %s", lastFileNum, err)
		return nil, err
	}

	if numBlocksInFile == 0 && lastFileNum > 0 {
		secondLastFileNum := lastFileNum - 1
		fileInfo := getFileInfoOrPanic(rootDir, secondLastFileNum)
		logger.Debugf("Second last Block file info: FileName=[%s], FileSize=[%d]", fileInfo.Name(), fileInfo.Size())
		if lastBlockBytes, _, _, err = scanForLastCompleteBlock(rootDir, secondLastFileNum, 0); err != nil {
			logger.Errorf("Error scanning second last file [num=%d]: %s", secondLastFileNum, err)
			return nil, err
		}
	}

	if lastBlockBytes != nil {
		if lastBlock, err = deserializeBlock(lastBlockBytes); err != nil {
			logger.Errorf("Error deserializing last block: %s. Block bytes length: %d", err, len(lastBlockBytes))
			return nil, err
		}
		lastBlockNumber = lastBlock.Header.Number
	}

	cpInfo := &checkpointInfo{
		lastBlockNumber:          lastBlockNumber,
		latestFileChunksize:      int(endOffsetLastBlock),
		latestFileChunkSuffixNum: lastFileNum,
		isChainEmpty:             lastFileNum == 0 && numBlocksInFile == 0,
	}
	logger.Debugf("Checkpoint info constructed from file system = %s", spew.Sdump(cpInfo))
	return cpInfo, nil
}

func syncArchiveMetaInfoFromDfs(rootDir string, conf *la.Config, amInfo *archive.ArchiveMetaInfo, client dc.FsClient) {
	logger.Infof("syncAMInfoFromDfs amInfo=%s", spew.Sdump(amInfo))
	//Checks if the file suffix of where the last block was written exists
	filePath := deriveBlockfilePath(rootDir, int(amInfo.LastSentFileSuffix))
	remotePath := conf.FsRoot + filePath
	if _, err := client.Stat(remotePath); err != nil {
		logger.Errorf("Error in checking whether file [%s] exists: %s", remotePath, err)
	} else {
		var lastSentFileNum int
		var lastArchiveFileNum int
		var err error
		if lastSentFileNum, err = retrieveLastSentFileNumFromDfs(rootDir, conf, client); err != nil {
			logger.Errorf("Error in retrieve last file suffix from dfs: %s", err)
		}
		if lastSentFileNum == int(amInfo.LastSentFileSuffix) {
			// archive meta info is in sync with the file on dfs
			return
		}

		//Scan the dfs to sync the archive meta info
		if lastArchiveFileNum, err = retrieveLastArchiveFileNum(rootDir); err != nil {
			logger.Errorf("Error in retrieve first file suffix from file system: %s", err)
		}
		amInfo.LastArchiveFileSuffix = int32(lastArchiveFileNum)
		amInfo.LastSentFileSuffix = int32(lastSentFileNum)
		calcFileProofs(rootDir, amInfo)
	}
}

func constructArchiveMetaInfoFromDfsBlockFiles(rootDir string, conf *la.Config, client dc.FsClient) (*archive.ArchiveMetaInfo, error) {
	logger.Info("Retrieving archive meta info from dfs block files")
	var lastArchiveFileNum int
	var lastSentFileNum int
	var err error
	amInfo := &archive.ArchiveMetaInfo{
		LastSentFileSuffix:    -1,
		LastArchiveFileSuffix: -1,
		FileProofs:            make(map[int32]string),
	}

	if lastArchiveFileNum, err = retrieveLastArchiveFileNum(rootDir); err != nil {
		return amInfo, err
	}
	if lastSentFileNum, err = retrieveLastSentFileNumFromDfs(rootDir, conf, client); err != nil {
		return amInfo, err
	}

	amInfo.LastArchiveFileSuffix = int32(lastArchiveFileNum)
	amInfo.LastSentFileSuffix = int32(lastSentFileNum)
	calcFileProofs(rootDir, amInfo)

	logger.Infof("Archive meta info constructed from dfs = %s", spew.Sdump(amInfo))
	return amInfo, nil
}

func retrieveLastArchiveFileNum(rootDir string) (int, error) {
	smallestFileNum := math.MaxInt32
	filesInfo, err := ioutil.ReadDir(rootDir)
	if err != nil {
		logger.Errorf("retrieveLastArchiveFileNum got error: %s", err)
		return -1, errors.Wrapf(err, "retrieveLastArchiveFileNum - error reading dir %s", rootDir)
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Infof("Skipping File name = %s", name)
			continue
		}
		fileNum, err := extractFileNum(name)
		if err != nil {
			return -1, err
		}
		if fileNum < smallestFileNum {
			smallestFileNum = fileNum
		}
	}
	if smallestFileNum == math.MaxInt32 {
		logger.Warnf("retrieveLastArchiveFileNum - no block file found in the dir[%s]", rootDir)
		return -1, nil
	}
	logger.Infof("retrieveLastArchiveFileNum() - smallestFileNum = %d", smallestFileNum)
	return smallestFileNum, err
}

func extractFileNum(name string) (int, error) {
	fileSuffix := strings.TrimPrefix(name, blockfilePrefix)
	fileNum, err := strconv.Atoi(fileSuffix)
	return fileNum, err
}

func calcFileProofs(rootDir string, amInfo *archive.ArchiveMetaInfo) {
	fileInfos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		logger.Errorf("ReadDir[%s] got err: %s", rootDir, err)
		return
	}

	for _, fi := range fileInfos {
		f, err := os.Open(rootDir + string(os.PathSeparator) + fi.Name())
		if err != nil {
			logger.Errorf("Open file[%s] got error: %s", fi.Name(), err)
			continue
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			logger.Errorf("io.Copy file[%s] got error: %s", fi.Name(), err)
			continue
		}
		fileNum, err := extractFileNum(fi.Name())
		if err != nil {
			logger.Errorf("extract file num from[%s] got error: %s", fi.Name(), err)
		}
		amInfo.FileProofs[int32(fileNum)] = hex.EncodeToString(h.Sum(nil))
	}
}

func retrieveLastSentFileNumFromDfs(rootDir string, conf *la.Config, client dc.FsClient) (int, error) {
	biggestFileNum := -1
	remotePath := conf.FsRoot + rootDir
	if _, notExistErr := client.Stat(remotePath); notExistErr != nil {
		logger.Warnf("retrieveLastSentFileNumFromDfs remotePath[%s] not exist", remotePath)
		return -1, notExistErr
	}
	filesInfo, err := client.ReadDir(remotePath)
	if err != nil {
		logger.Errorf("retrieveLastSentFileNumFromDfs got error: %s", err)
		return -1, errors.Wrapf(err, "retrieveLastSentFileNumFromDfs - error reading dir %s", remotePath)
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Infof("Skipping File name = %s", name)
			continue
		}
		fileNum, err := extractFileNum(name)
		if err != nil {
			return -1, err
		}
		if fileNum > biggestFileNum {
			biggestFileNum = fileNum
		}
	}
	if biggestFileNum == -1 {
		logger.Warnf("retrieveLastSentFileNumFromDfs - no block file found in the dfs dir[%s]", remotePath)
	}
	logger.Infof("retrieveLastSentFileNumFromDfs() - biggestFileNum = %d", biggestFileNum)
	return biggestFileNum, err
}

// binarySearchFileNumForBlock locates the file number that contains the given block number.
// This function assumes that the caller invokes this function with a block number that has been commited
// For any uncommitted block, this function returns the last file present
func binarySearchFileNumForBlock(rootDir string, blockNum uint64) (int, error) {
	cpInfo, err := constructCheckpointInfoFromBlockFiles(rootDir)
	if err != nil {
		return -1, err
	}

	beginFile := 0
	endFile := cpInfo.latestFileChunkSuffixNum

	for endFile != beginFile {
		searchFile := beginFile + (endFile-beginFile)/2 + 1
		n, err := retrieveFirstBlockNumFromFile(rootDir, searchFile)
		if err != nil {
			return -1, err
		}
		switch {
		case n == blockNum:
			return searchFile, nil
		case n > blockNum:
			endFile = searchFile - 1
		case n < blockNum:
			beginFile = searchFile
		}
	}
	return beginFile, nil
}

func retrieveFirstBlockNumFromFile(rootDir string, fileNum int) (uint64, error) {
	s, err := newBlockfileStream(rootDir, fileNum, 0)
	if err != nil {
		return 0, err
	}
	defer s.close()
	bb, err := s.nextBlockBytes()
	if err != nil {
		return 0, err
	}
	blockInfo, err := extractSerializedBlockInfo(bb)
	if err != nil {
		return 0, err
	}
	return blockInfo.blockHeader.Number, nil
}

func retrieveLastFileSuffix(rootDir string) (int, error) {
	logger.Debugf("retrieveLastFileSuffix()")
	biggestFileNum := -1
	filesInfo, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return -1, errors.Wrapf(err, "error reading dir %s", rootDir)
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Debugf("Skipping File name = %s", name)
			continue
		}
		fileNum, err := extractFileNum(name)
		if err != nil {
			return -1, err
		}
		if fileNum > biggestFileNum {
			biggestFileNum = fileNum
		}
	}
	logger.Debugf("retrieveLastFileSuffix() - biggestFileNum = %d", biggestFileNum)
	return biggestFileNum, err
}

func isBlockFileName(name string) bool {
	return strings.HasPrefix(name, blockfilePrefix)
}

func getFileInfoOrPanic(rootDir string, fileNum int) os.FileInfo {
	filePath := deriveBlockfilePath(rootDir, fileNum)
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		panic(errors.Wrapf(err, "error retrieving file info for file number %d", fileNum))
	}
	return fileInfo
}
