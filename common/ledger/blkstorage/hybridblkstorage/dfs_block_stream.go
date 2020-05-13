/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"bufio"
	"fmt"
	"github.com/pengisgood/hdfs"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// ErrUnexpectedEndOfDfsBlockfile error used to indicate an unexpected end of a file segment
// this can happen mainly if a crash occurs during appening a block and partial block contents
// get written towards the end of the file
var ErrUnexpectedEndOfDfsBlockfile = errors.New("unexpected end of dfs blockfile")

type fileStream interface {
	nextBlockBytes() ([]byte, error)
	close() error
}

// dfsBlockfileStream reads blocks sequentially from a single file.
// It starts from the given offset and can traverse till the end of the file
type dfsBlockfileStream struct {
	fileNum       int
	reader        *hdfs.FileReader
	currentOffset int64
}

///////////////////////////////////
// dfsBlockfileStream functions
////////////////////////////////////
func newDfsBlockfileStream(rootDir string, fileNum int, startOffset int64, client *hdfs.Client) (*dfsBlockfileStream, error) {
	var reader *hdfs.FileReader
	var err error

	filePath := deriveBlockfilePath(rootDir, fileNum)
	logger.Infof("newDfsBlockfileStream(): filePath=[%s], startOffset=[%d]", filePath, startOffset)
	if reader, err = client.Open(filePath); err != nil {
		return nil, errors.Wrapf(err, "error opening dfs block reader %s", filePath)
	}
	var newPosition int64
	if newPosition, err = reader.Seek(startOffset, 0); err != nil {
		return nil, errors.Wrapf(err, "error seeking dfs block reader [%s] to startOffset [%d]", filePath, startOffset)
	}
	if newPosition != startOffset {
		panic(fmt.Sprintf("Could not seek dfs block reader [%s] to startOffset [%d]. New position = [%d]",
			filePath, startOffset, newPosition))
	}

	s := &dfsBlockfileStream{fileNum, reader, startOffset}
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
