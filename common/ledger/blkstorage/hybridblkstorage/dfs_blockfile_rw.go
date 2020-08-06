/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/pkg/errors"
	"io"
)

type fileReader interface {
	readAt(offset int, length int) ([]byte, error)
	close() error
}

////  DFS READER ////
type dfsBlockfileReader struct {
	reader common.FsReader
}

func newDfsBlockfileReader(filePath string, dfsClient common.FsClient, fileProof string) (*dfsBlockfileReader, error) {
	if dfsClient == nil {
		return nil, errors.New("dfs client should not be nil")
	}
	reader, err := dfsClient.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening dfs block file reader for file %s with client %+v", filePath, dfsClient)
	}

	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		logger.Errorf("io.Copy file[%s] got error: %s", filePath, err)
	} else if len(fileProof) != 0 && hex.EncodeToString(h.Sum(nil)) != fileProof {
		logger.Warnf("the checksum of blockfile[%s] in dfs does not match [%s]", filePath, fileProof)
	}

	//reset the read position to the beginning
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		logger.Errorf("newDfsBlockfileReader - reset read position of file[%s] got error: %s", filePath, err)
		return nil, errors.Wrapf(err, "newDfsBlockfileReader - reset read position of file[%s]", filePath)
	}

	return &dfsBlockfileReader{reader}, nil
}

func (r *dfsBlockfileReader) readAt(offset int, length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := r.reader.ReadAt(b, int64(offset))
	if err != nil {
		return nil, errors.Wrapf(err, "error reading dfs block file for offset %d and length %d", offset, length)
	}
	return b, nil
}

func (r *dfsBlockfileReader) close() error {
	return errors.WithStack(r.reader.Close())
}
