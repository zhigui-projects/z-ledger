/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"github.com/pengisgood/hdfs"
	"github.com/pkg/errors"
)

type FileReader interface {
	ReadAt(offset int, length int) ([]byte, error)
	Close() error
}

////  DFS READER ////
type dfsBlockfileReader struct {
	reader *hdfs.FileReader
}

func newDfsBlockfileReader(filePath string, dfsClient *hdfs.Client) (*dfsBlockfileReader, error) {
	if dfsClient == nil {
		return nil, errors.New("dfs client should not be nil")
	}
	reader, err := dfsClient.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening dfs block file reader for file %s with client %+v", filePath, dfsClient)
	}
	return &dfsBlockfileReader{reader}, nil
}

func (r *dfsBlockfileReader) ReadAt(offset int, length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := r.reader.ReadAt(b, int64(offset))
	if err != nil {
		return nil, errors.Wrapf(err, "error reading dfs block file for offset %d and length %d", offset, length)
	}
	return b, nil
}

func (r *dfsBlockfileReader) Close() error {
	return errors.WithStack(r.reader.Close())
}
