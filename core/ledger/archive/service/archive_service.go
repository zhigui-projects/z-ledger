/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"fmt"
	"github.com/colinmarc/hdfs"
	"github.com/fsnotify/fsnotify"
	proto "github.com/hyperledger/fabric-protos-go/gossip"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger/blkstorage/hybridblkstorage"
	coreconfig "github.com/hyperledger/fabric/core/config"
	le "github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/archive/dfs"
	"github.com/hyperledger/fabric/core/ledger/kvledger"
	"github.com/hyperledger/fabric/core/ledger/ledgermgmt"
	"github.com/pkg/errors"
	"path/filepath"
	"sync"
)

var logger = flogging.MustGetLogger("archive.service")

type gossip interface {
	// PeersOfChannel returns the NetworkMembers considered alive in a channel
	//PeersOfChannel(channel common.ChannelID) []discovery.NetworkMember

	// Accept returns a dedicated read-only channel for messages sent by other nodes that match a certain predicate.
	// If passThrough is false, the messages are processed by the gossip layer beforehand.
	// If passThrough is true, the gossip layer doesn't intervene and the messages
	// can be used to send a reply back to the sender
	//Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan protoext.ReceivedMessage)

	// Gossip sends a message to other peers to the network
	Gossip(msg *proto.GossipMessage)
}

type archiveSvc interface {
	// StartWatcherForChannel dynamically starts watcher of ledger files.
	StartWatcherForChannel(chainID string) error

	// StopWatcherForChannel dynamically stops watcher of ledger files.
	StopWatcherForChannel(chainID string) error

	// Stop stops watcher of ledger files.
	Stop()
}

type ArchiveService struct {
	gossip
	archiveSvc
	ledgerMgr *ledgermgmt.LedgerMgr
	dfsClient *hdfs.Client
	watchers  map[string]*fsnotify.Watcher
	lock      sync.RWMutex
}

// New construction function to create and initialize
// archive service instance. It tries to establish connection to
// the specified dfs name node, in case it fails to dial to it, return nil
func New(g gossip, ledgerMgr *ledgermgmt.LedgerMgr) (*ArchiveService, error) {
	client, err := dfs.NewHDFSClient()
	if err != nil {
		logger.Errorf("Archive service can't connect to HDFS, due to %+v", err)
		return nil, err
	}

	return &ArchiveService{
		gossip:    g,
		ledgerMgr: ledgerMgr,
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

	go a.broadcastMetaInfo(chainID)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Errorf("Archive service - watcher has been closed")
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					logger.Infof("Archive service - new ledger file CREATED event: %s", event.Name)
					go a.transferBlockFiles(chainID)
					go a.broadcastMetaInfo(chainID)
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

func (a *ArchiveService) transferBlockFiles(chainID string) {
	ledger, err := a.getLedger(chainID)
	if err != nil {
		logger.Errorf("Archive service - get ledger for chain: %s failed with error: %s", chainID, err)
	}

	err = ledger.TransferBlockFiles()
	if err != nil {
		logger.Errorf("Archive service - transfer block files for chain: %s failed with error: %s", chainID, err)
	}
}

func (a *ArchiveService) broadcastMetaInfo(chainID string) {
	ledger, err := a.getLedger(chainID)
	if err != nil {
		logger.Errorf("Archive service - get ledger for chain: %s failed with error: %s", chainID, err)
		return
	}

	metaInfo, err := ledger.GetArchiveMetaInfo()
	if err != nil {
		logger.Errorf("Archive service - get meta info failed with error: %s", err)
		return
	}
	if metaInfo == nil {
		logger.Error("Archive service - meta info initialization failed")
		return
	}

	// Use payload to create gossip message
	gossipMsg := &proto.GossipMessage{
		Nonce:   0,
		Tag:     proto.GossipMessage_CHAN_ONLY,
		Channel: []byte(chainID),
		Content: &proto.GossipMessage_ArchiveMsg{
			ArchiveMsg: &proto.ArchiveMsg{
				ArchiveMetaInfo: metaInfo,
			},
		},
	}

	// Gossip messages with other nodes
	logger.Infof("Archive service - gossiping archive meta info: %+v", metaInfo)
	a.gossip.Gossip(gossipMsg)
}

func (a *ArchiveService) getLedger(chainID string) (le.PeerLedger, error) {
	ledger, err := a.ledgerMgr.GetOpenedLedger(chainID)
	if err != nil {
		logger.Warnf("Archive service - get opened ledger for chain: %s failed with error: %s", chainID, err)

		ledger, err = a.ledgerMgr.OpenLedger(chainID)
		if err != nil {
			logger.Errorf("Archive service - open ledger for chain: %s failed with error: %s", chainID, err)
		}
	}
	return ledger, err
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
