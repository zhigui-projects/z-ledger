/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"path/filepath"
)

const (
	// ChainsDir is the name of the directory containing the channel ledgers.
	ChainsDir = "chains"
	// IndexDir is the name of the directory containing all block indexes across ledgers.
	IndexDir = "index"
)

// Conf encapsulates all the configurations for `FsBlockStore`
type Conf struct {
	blockStorageDir  string
	maxBlockfileSize int
}

// NewConf constructs new `Conf`.
// blockStorageDir is the top level folder under which `HybridBlockStore` manages its data
func NewConf(blockStorageDir string, maxBlockfileSize int) *Conf {
	if maxBlockfileSize <= 0 {
		maxBlockfileSize = ledgerconfig.GetDefaultMaxBlockfileSize()
	}
	return &Conf{blockStorageDir, maxBlockfileSize}
}

func (conf *Conf) getIndexDir() string {
	return filepath.Join(conf.blockStorageDir, IndexDir)
}

func (conf *Conf) getChainsDir() string {
	return filepath.Join(conf.blockStorageDir, ChainsDir)
}

func (conf *Conf) getLedgerBlockDir(ledgerid string) string {
	return filepath.Join(conf.getChainsDir(), ledgerid)
}
