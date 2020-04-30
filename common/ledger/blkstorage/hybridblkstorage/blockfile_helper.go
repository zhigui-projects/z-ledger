/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/colinmarc/hdfs"
	"github.com/davecgh/go-spew/spew"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/ledger/blkstorage/hybridblkstorage/msgs"
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

func syncArchiveMetaInfoFromDfs(rootDir string, amInfo *msgs.ArchiveMetaInfo, client *hdfs.Client) {
	logger.Infof("syncAMInfoFromDfs amInfo=%s", spew.Sdump(amInfo))
	//Checks if the file suffix of where the last block was written exists
	filePath := deriveBlockfilePath(rootDir, int(amInfo.LastSentFileSuffix))
	if _, err := client.Stat(filePath); err != nil {
		panic(fmt.Sprintf("Error in checking whether file [%s] exists: %s", filePath, err))
	}

	var lastSentFileNum int
	var lastArchiveFileNum int
	var err error
	if lastSentFileNum, err = retrieveLastFileSuffixFromDfs(rootDir, client); err != nil {
		panic(fmt.Sprintf("Error in retrieve last file suffix from dfs: %s", err))
	}
	if lastSentFileNum == int(amInfo.LastSentFileSuffix) {
		// archive meta info is in sync with the file on dfs
		return
	}

	//Scan the dfs to sync the archive meta info
	if lastArchiveFileNum, err = retrieveFirstFileSuffix(rootDir); err != nil {
		panic(fmt.Sprintf("Error in retrieve first file suffix from file system: %s", err))
	}
	amInfo.LastArchiveFileSuffix = int32(lastArchiveFileNum - 1)
	amInfo.LastSentFileSuffix = int32(lastSentFileNum)
	//TODO：calculate checksum
	//amInfo.FileProofs[]
}

func constructArchiveMetaInfoFromDfsBlockFiles(rootDir string, client *hdfs.Client) (*msgs.ArchiveMetaInfo, error) {
	logger.Info("Retrieving archive meta info from dfs block files")
	var lastArchiveFileNum int
	var lastSentFileNum int
	var err error
	amInfo := &msgs.ArchiveMetaInfo{
		LastSentFileSuffix:    -1,
		LastArchiveFileSuffix: -1,
		FileProofs:            make(map[int32]string),
	}

	if lastArchiveFileNum, err = retrieveFirstFileSuffix(rootDir); err != nil {
		return nil, err
	}
	if lastSentFileNum, err = retrieveLastFileSuffixFromDfs(rootDir, client); err != nil {
		return nil, err
	}

	if lastArchiveFileNum == math.MaxInt32 || lastSentFileNum == -1 {
		logger.Warn("No block file found")
		return amInfo, nil
	}

	amInfo.LastArchiveFileSuffix = int32(lastArchiveFileNum - 1)
	amInfo.LastSentFileSuffix = int32(lastSentFileNum)
	//TODO：calculate checksum
	//amInfo.FileProofs[]
	logger.Infof("Archive meta info constructed from dfs = %s", spew.Sdump(amInfo))
	return amInfo, nil
}

func retrieveFirstFileSuffix(rootDir string) (int, error) {
	smallestFileNum := math.MaxInt32
	filesInfo, err := ioutil.ReadDir(rootDir)
	if err != nil {
		logger.Errorf("retrieveFirstFileSuffix got error: %s", err)
		return -1, errors.Wrapf(err, "retrieveFirstFileSuffix - error reading dir %s", rootDir)
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Infof("Skipping File name = %s", name)
			continue
		}
		fileSuffix := strings.TrimPrefix(name, blockfilePrefix)
		fileNum, err := strconv.Atoi(fileSuffix)
		if err != nil {
			return -1, err
		}
		if fileNum < smallestFileNum {
			smallestFileNum = fileNum
		}
	}
	logger.Infof("retrieveFirstFileSuffix() - smallestFileNum = %d", smallestFileNum)
	return smallestFileNum, err
}

func retrieveLastFileSuffixFromDfs(rootDir string, client *hdfs.Client) (int, error) {
	biggestFileNum := -1
	filesInfo, err := client.ReadDir(rootDir)
	if err != nil {
		logger.Errorf("retrieveLastFileSuffixFromDfs got error: %s", err)
		return -1, errors.Wrapf(err, "retrieveLastFileSuffixFromDfs - error reading dir %s", rootDir)
	}
	for _, fileInfo := range filesInfo {
		name := fileInfo.Name()
		if fileInfo.IsDir() || !isBlockFileName(name) {
			logger.Infof("Skipping File name = %s", name)
			continue
		}
		fileSuffix := strings.TrimPrefix(name, blockfilePrefix)
		fileNum, err := strconv.Atoi(fileSuffix)
		if err != nil {
			return -1, err
		}
		if fileNum > biggestFileNum {
			biggestFileNum = fileNum
		}
	}
	logger.Infof("retrieveLastFileSuffixFromDfs() - biggestFileNum = %d", biggestFileNum)
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
		fileSuffix := strings.TrimPrefix(name, blockfilePrefix)
		fileNum, err := strconv.Atoi(fileSuffix)
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
