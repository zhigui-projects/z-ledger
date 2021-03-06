/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hybridblkstorage

import (
	"os"

	"github.com/hyperledger/fabric/common/ledger/blkstorage"
	"github.com/hyperledger/fabric/common/ledger/dataformat"
	"github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/archive"
	"github.com/pkg/errors"
)

func dataFormatVersion(indexConfig *blkstorage.IndexConfig) string {
	// in version 2.0 we merged three indexable into one `IndexableAttrTxID`
	if indexConfig.Contains(blkstorage.IndexableAttrTxID) {
		return dataformat.Version20
	}
	return dataformat.Version1x
}

// HybridBlockstoreProvider provides handle to block storage - this is not thread-safe
type HybridBlockstoreProvider struct {
	conf            *Conf
	indexConfig     *blkstorage.IndexConfig
	leveldbProvider *leveldbhelper.Provider
	stats           *stats
	archiveConf     *archive.Config
}

// NewProvider constructs a filesystem based block store provider
func NewProvider(conf *Conf, indexConfig *blkstorage.IndexConfig, metricsProvider metrics.Provider,
	archiveConf *archive.Config) (blkstorage.BlockStoreProvider, error) {
	dbConf := &leveldbhelper.Conf{
		DBPath:                conf.getIndexDir(),
		ExpectedFormatVersion: dataFormatVersion(indexConfig),
	}

	p, err := leveldbhelper.NewProvider(dbConf)
	if err != nil {
		return nil, err
	}

	dirPath := conf.getChainsDir()
	if _, err := os.Stat(dirPath); err != nil {
		if !os.IsNotExist(err) { // NotExist is the only permitted error type
			return nil, errors.Wrapf(err, "failed to read ledger directory %s", dirPath)
		}

		logger.Info("Creating new file ledger directory at", dirPath)
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return nil, errors.Wrapf(err, "failed to create ledger directory: %s", dirPath)
		}
	}

	// create stats instance at provider level and pass to newHybridBlockStore
	stats := newStats(metricsProvider)
	return &HybridBlockstoreProvider{conf, indexConfig, p, stats, archiveConf}, nil
}

// CreateBlockStore simply calls OpenBlockStore
func (p *HybridBlockstoreProvider) CreateBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	return p.OpenBlockStore(ledgerid)
}

// OpenBlockStore opens a block store for given ledgerid.
// If a blockstore is not existing, this method creates one
// This method should be invoked only once for a particular ledgerid
func (p *HybridBlockstoreProvider) OpenBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	indexStoreHandle := p.leveldbProvider.GetDBHandle(ledgerid)
	return newHybridBlockStore(ledgerid, p.conf, p.indexConfig, indexStoreHandle, p.stats, p.archiveConf), nil
}

// Exists tells whether the BlockStore with given id exists
func (p *HybridBlockstoreProvider) Exists(ledgerid string) (bool, error) {
	exists, _, err := util.FileExists(p.conf.getLedgerBlockDir(ledgerid))
	return exists, err
}

// List lists the ids of the existing ledgers
func (p *HybridBlockstoreProvider) List() ([]string, error) {
	return util.ListSubdirs(p.conf.getChainsDir())
}

// Close closes the HybridBlockstoreProvider
func (p *HybridBlockstoreProvider) Close() {
	p.leveldbProvider.Close()
}
