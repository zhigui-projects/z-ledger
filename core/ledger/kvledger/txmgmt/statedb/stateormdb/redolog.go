/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"bytes"
	"encoding/gob"

	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
)

var redologKeyPrefix = []byte{byte(0)}

type RedoLoggerProvider struct {
	leveldbProvider *leveldbhelper.Provider
}
type RedoLogger struct {
	dbName   string
	dbHandle *leveldbhelper.DBHandle
}

type RedoRecord struct {
	UpdateBatch *statedb.UpdateBatch
	Version     *version.Height
}

func NewRedoLoggerProvider(dirPath string) (*RedoLoggerProvider, error) {
	provider, err := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dirPath})
	if err != nil {
		return nil, err
	}
	return &RedoLoggerProvider{leveldbProvider: provider}, nil
}

func (p *RedoLoggerProvider) NewRedoLogger(dbName string) *RedoLogger {
	return &RedoLogger{
		dbHandle: p.leveldbProvider.GetDBHandle(dbName),
	}
}

func (p *RedoLoggerProvider) Close() {
	p.leveldbProvider.Close()
}

func (l *RedoLogger) Persist(r *RedoRecord) error {
	k := encodeRedologKey(l.dbName)
	v, err := EncodeRedologVal(r)
	if err != nil {
		return err
	}
	return l.dbHandle.Put(k, v, true)
}

func (l *RedoLogger) Load() (*RedoRecord, error) {
	k := encodeRedologKey(l.dbName)
	v, err := l.dbHandle.Get(k)
	if err != nil || v == nil {
		return nil, err
	}
	return decodeRedologVal(v)
}

func encodeRedologKey(dbName string) []byte {
	return append(redologKeyPrefix, []byte(dbName)...)
}

func EncodeRedologVal(r *RedoRecord) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeRedologVal(b []byte) (*RedoRecord, error) {
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	var r *RedoRecord
	if err := decoder.Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
