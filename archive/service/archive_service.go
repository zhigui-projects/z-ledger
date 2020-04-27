/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"fmt"
	"github.com/colinmarc/hdfs"
	"github.com/fsnotify/fsnotify"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger/blkstorage/hybridblkstorage"
	coreconfig "github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/core/ledger/kvledger"
	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"github.com/pkg/errors"
	"path/filepath"
	"sync"
)

var logger = flogging.MustGetLogger("archive.service")

type archiveSvc interface {
	// StartWatcherForChannel dynamically starts watcher of ledger files.
	StartWatcherForChannel(chainID string) error

	// StopWatcherForChannel dynamically stops watcher of ledger files.
	StopWatcherForChannel(chainID string) error

	// Stop stops watcher of ledger files.
	Stop()
}

type ArchiveService struct {
	archiveSvc
	dfsClient *hdfs.Client
	watchers  map[string]*fsnotify.Watcher
	lock      sync.RWMutex
}

// New construction function to create and initialize
// archive service instance. It tries to establish connection to
// the specified dfs name node, in case it fails to dial to it, return nil
func New() (*ArchiveService, error) {
	if len(ledgerconfig.GetHDFSNameNodes()) == 0 {
		errMsg := "Archive service can't be initialized, due to no namenode address in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if len(ledgerconfig.GetHDFSUser()) == 0 {
		errMsg := "Archive service can't be initialized, due to no HDFS user in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	client, err := hdfs.NewClient(hdfs.ClientOptions{
		Addresses: ledgerconfig.GetHDFSNameNodes(),
		User:      ledgerconfig.GetHDFSUser(),
	})

	if err != nil {
		logger.Errorf("Archive service can't connect to HDFS, due to %+v", err)
		return nil, err
	}

	return &ArchiveService{
		dfsClient: client,
		watchers:  make(map[string]*fsnotify.Watcher),
	}, nil
}

// StartWatcherForChannel starts ledger watcher for channel
func (a *ArchiveService) StartWatcherForChannel(chainID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, exist := a.watchers[chainID]; exist {
		errMsg := fmt.Sprintf("Archive service - ledger watcher already exists for %s found, can't start a new watcher", chainID)
		logger.Warn(errMsg)
		return errors.New(errMsg)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Errorf("Archive service - can't start watcher, due to %+v", err)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Errorf("Archive service - watcher has been closed")
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					logger.Infof("Archive service - created ledger file: %s", event.Name)
					//TODO: transfer file
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logger.Errorf("Archive service - watcher has been closed, due to err: %+v", err)
					return
				}
				logger.Errorf("Archive service - watcher got error: %+v", err)
			}
		}
	}()

	ledgerDir := getLedgerDir(chainID)
	logger.Infof("Archive service - adding watcher for ledger directory: %s", ledgerDir)
	err = watcher.Add(ledgerDir)
	if err != nil {
		errMsg := fmt.Sprintf("Archive service - add watching directory failed, due to %+v", err)
		logger.Error(errMsg)
		return errors.New(errMsg)
	}
	a.watchers[chainID] = watcher

	return nil
}

func (a *ArchiveService) StopWatcherForChannel(chainID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	watcher, exist := a.watchers[chainID]
	if !exist {
		errMsg := fmt.Sprintf("Archive service - ledger watcher does not exists for %s found", chainID)
		logger.Warn(errMsg)
		return errors.New(errMsg)
	}

	if err := watcher.Close(); err != nil {
		errMsg := fmt.Sprintf("Archive service - close watcher failed for channel: %s, due to %+v", chainID, err)
		logger.Error(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (a *ArchiveService) Stop() error {
	//TODO
	return nil
}

func getLedgerDir(chainID string) string {
	rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
	chainsDir := filepath.Join(kvledger.BlockStorePath(rootFSPath), hybridblkstorage.ChainsDir)
	ledgerDir := filepath.Join(chainsDir, chainID)
	return ledgerDir
}
