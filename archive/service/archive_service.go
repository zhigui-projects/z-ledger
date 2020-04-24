/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"github.com/colinmarc/hdfs"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"github.com/pkg/errors"
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

	return &ArchiveService{dfsClient: client}, nil
}

// StartWatcherForChannel starts ledger watcher for channel
func (a *ArchiveService) StartWatcherForChannel(chainID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	//TODO
	return nil
}

func (a *ArchiveService) StopWatcherForChannel(chainID string) error {
	//TODO
	return nil
}

func (a *ArchiveService) Stop() error {
	//TODO
	return nil
}
