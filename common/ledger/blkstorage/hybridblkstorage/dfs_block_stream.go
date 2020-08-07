/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/pkg/errors"
)

// ErrUnexpectedEndOfDfsBlockfile error used to indicate an unexpected end of a file segment
// this can happen mainly if a crash occurs during appening a block and partial block contents
// get written towards the end of the file
var ErrUnexpectedEndOfDfsBlockfile = errors.New("unexpected end of dfs blockfile")

type hybridBlockStream interface {
	nextBlockBytes() ([]byte, error)
	close() error
}

// dfsBlockfileStream reads blocks sequentially from a single file.
// It starts from the given offset and can traverse till the end of the file
type dfsBlockfileStream struct {
	fileNum       int
	reader        common.FsReader
	client        common.FsClient
	currentOffset int64
}

// dfsBlockStream reads blocks sequentially from multiple files.
// it starts from a given file offset and continues with the next
// file segment until the end of the last segment (`endFileNum`)
type dfsBlockStream struct {
	rootDir           string
	currentFileNum    int
	endFileNum        int
	currentFileStream *dfsBlockfileStream
}

///////////////////////////////////
// dfsBlockfileStream functions
////////////////////////////////////
func newDfsBlockfileStream(rootDir string, fileNum int, startOffset int64, client common.FsClient, fileProof string) (*dfsBlockfileStream, error) {
	var reader common.FsReader
	var err error

	filePath := deriveBlockfilePath(rootDir, fileNum)
	logger.Infof("newDfsBlockfileStream(): filePath=[%s], startOffset=[%d]", filePath, startOffset)
	if reader, err = client.Open(filePath); err != nil {
		return nil, errors.Wrapf(err, "error opening dfs block reader %s", filePath)
	}
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		logger.Errorf("io.Copy file[%s] got error: %s", filePath, err)
	} else if len(fileProof) != 0 && hex.EncodeToString(h.Sum(nil)) != fileProof {
		logger.Warnf("the checksum of blockfile[%s] in dfs does not match [%s]", filePath, fileProof)
	}

	//reset the read position to the beginning
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		logger.Errorf("newDfsBlockfileStream - reset read position of file[%s] got error: %s", filePath, err)
		return nil, errors.Wrapf(err, "newDfsBlockfileStream - reset read position of file[%s]", filePath)
	}

	var newPosition int64
	if newPosition, err = reader.Seek(startOffset, io.SeekStart); err != nil {
		logger.Errorf("reader.Seek file[%s], offset[%d] in dfs got error: %+v", filePath, startOffset, err)
		return nil, errors.Wrapf(err, "error seeking dfs block reader [%s] to startOffset [%d]", filePath, startOffset)
	}
	if newPosition != startOffset {
		panic(fmt.Sprintf("Could not seek dfs block reader [%s] to startOffset [%d]. New position = [%d]",
			filePath, startOffset, newPosition))
	}

	s := &dfsBlockfileStream{fileNum, reader, client, startOffset}
	return s, nil
}

func (s *dfsBlockfileStream) nextBlockBytes() ([]byte, error) {
	blockBytes, _, err := s.nextBlockBytesAndPlacementInfo()
	return blockBytes, err
}

// nextBlockBytesAndPlacementInfo returns bytes for the next block
// along with the offset information in the block file.
// An error `ErrUnexpectedEndOfDfsBlockfile` is returned if a partial written data is detected
// which is possible towards the tail of the file if a crash had taken place during appending of a block
func (s *dfsBlockfileStream) nextBlockBytesAndPlacementInfo() ([]byte, *blockPlacementInfo, error) {
	var lenBytes []byte
	var err error
	var fileInfo os.FileInfo

	moreContentAvailable := true
	fileInfo = s.reader.Stat()
	if s.currentOffset == fileInfo.Size() {
		logger.Infof("Finished reading dfs file number [%d]", s.fileNum)
		return nil, nil, nil
	}
	remainingBytes := fileInfo.Size() - s.currentOffset
	// Peek 8 or smaller number of bytes (if remaining bytes are less than 8)
	// Assumption is that a block size would be small enough to be represented in 8 bytes varint
	peekBytes := 8
	if remainingBytes < int64(peekBytes) {
		peekBytes = int(remainingBytes)
		moreContentAvailable = false
	}
	logger.Infof("Remaining bytes=[%d], Going to peek [%d] bytes", remainingBytes, peekBytes)
	reader := bufio.NewReader(s.reader)
	if lenBytes, err = reader.Peek(peekBytes); err != nil {
		return nil, nil, errors.Wrapf(err, "error peeking [%d] bytes from dfs block file", peekBytes)
	}
	length, n := proto.DecodeVarint(lenBytes)
	if n == 0 {
		// proto.DecodeVarint did not consume any byte at all which means that the bytes
		// representing the size of the block are partial bytes
		if !moreContentAvailable {
			return nil, nil, ErrUnexpectedEndOfDfsBlockfile
		}
		panic(errors.Errorf("Error in decoding varint bytes [%#v]", lenBytes))
	}
	bytesExpected := int64(n) + int64(length)
	if bytesExpected > remainingBytes {
		logger.Debugf("At least [%d] bytes expected. Remaining bytes = [%d]. Returning with error [%s]",
			bytesExpected, remainingBytes, ErrUnexpectedEndOfBlockfile)
		return nil, nil, ErrUnexpectedEndOfDfsBlockfile
	}
	// skip the bytes representing the block size
	if _, err = reader.Discard(n); err != nil {
		return nil, nil, errors.Wrapf(err, "error discarding [%d] bytes", n)
	}
	blockBytes := make([]byte, length)
	if _, err = io.ReadAtLeast(reader, blockBytes, int(length)); err != nil {
		logger.Errorf("Error reading [%d] bytes from dfs block file number [%d], error: %s", length, s.fileNum, err)
		return nil, nil, errors.Wrapf(err, "error reading [%d] bytes from dfs block file number [%d]", length, s.fileNum)
	}
	blockPlacementInfo := &blockPlacementInfo{
		fileNum:          s.fileNum,
		blockStartOffset: s.currentOffset,
		blockBytesOffset: s.currentOffset + int64(n)}
	s.currentOffset += int64(n) + int64(length)
	logger.Infof("Returning blockbytes - length=[%d], placementInfo={%s}", len(blockBytes), blockPlacementInfo)
	return blockBytes, blockPlacementInfo, nil
}

func (s *dfsBlockfileStream) close() error {
	return errors.WithStack(s.reader.Close())
}

///////////////////////////////////
// dfsBlockStream functions
////////////////////////////////////
func newDfsBlockStream(rootDir string, startFileNum int, startOffset int64, endFileNum int, client common.FsClient) (*dfsBlockStream, error) {
	startFileStream, err := newDfsBlockfileStream(rootDir, startFileNum, startOffset, client, "")
	if err != nil {
		return nil, err
	}
	return &dfsBlockStream{rootDir, startFileNum, endFileNum, startFileStream}, nil
}

func (s *dfsBlockStream) moveToNextBlockfileStream() error {
	var err error
	if err = s.currentFileStream.close(); err != nil {
		return err
	}
	s.currentFileNum++
	if s.currentFileStream, err = newDfsBlockfileStream(s.rootDir, s.currentFileNum, 0, s.currentFileStream.client, ""); err != nil {
		return err
	}
	return nil
}

func (s *dfsBlockStream) nextBlockBytes() ([]byte, error) {
	blockBytes, _, err := s.nextBlockBytesAndPlacementInfo()
	return blockBytes, err
}

func (s *dfsBlockStream) nextBlockBytesAndPlacementInfo() ([]byte, *blockPlacementInfo, error) {
	var blockBytes []byte
	var blockPlacementInfo *blockPlacementInfo
	var err error
	if blockBytes, blockPlacementInfo, err = s.currentFileStream.nextBlockBytesAndPlacementInfo(); err != nil {
		logger.Errorf("Error reading next dfs block bytes from file number [%d]: %s", s.currentFileNum, err)
		return nil, nil, err
	}
	logger.Debugf("blockbytes [%d] read from dfs block file [%d]", len(blockBytes), s.currentFileNum)
	if blockBytes == nil && (s.currentFileNum < s.endFileNum || s.endFileNum < 0) {
		logger.Debugf("current file [%d] exhausted. Moving to next dfs block file", s.currentFileNum)
		if err = s.moveToNextBlockfileStream(); err != nil {
			return nil, nil, err
		}
		return s.nextBlockBytesAndPlacementInfo()
	}
	return blockBytes, blockPlacementInfo, nil
}

func (s *dfsBlockStream) close() error {
	return s.currentFileStream.close()
}
