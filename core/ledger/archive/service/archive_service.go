/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	proto "github.com/hyperledger/fabric-protos-go/gossip"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger/blkstorage/hybridblkstorage"
	coreconfig "github.com/hyperledger/fabric/core/config"
	le "github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/archive/dfs"
	"github.com/hyperledger/fabric/core/ledger/kvledger"
	"github.com/hyperledger/fabric/core/ledger/ledgermgmt"
	"github.com/pengisgood/hdfs"
	"github.com/pkg/errors"
	"path/filepath"
	"sync"
)

var logger = flogging.MustGetLogger("archive.service")

type ArchiveGossip interface {
	// Accept returns a channel that emits messages
	Accept() <-chan *proto.GossipMessage

	// Gossip gossips a message to other peers
	Gossip(msg *proto.GossipMessage)

	Stop()
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
	archiveSvc
	gossip    ArchiveGossip
	peerId    string
	ledgerMgr *ledgermgmt.LedgerMgr
	dfsClient *hdfs.Client
	watchers  map[string]*fsnotify.Watcher
	lock      sync.RWMutex
	stopCh    chan struct{}
}

// New construction function to create and initialize
// archive service instance. It tries to establish connection to
// the specified dfs name node, in case it fails to dial to it, return nil
func New(g gossip, ledgerMgr *ledgermgmt.LedgerMgr, channel string, peerId string) (*ArchiveService, error) {
	client, err := dfs.NewHDFSClient()
	if err != nil {
		logger.Errorf("Archive service can't connect to HDFS, due to %+v", err)
		return nil, err
	}

	a := &ArchiveService{
		gossip:    NewGossipImpl(g, channel),
		peerId:    peerId,
		ledgerMgr: ledgerMgr,
		dfsClient: client,
		watchers:  make(map[string]*fsnotify.Watcher),
		stopCh:    make(chan struct{}),
	}
	go a.handleMessages()

	return a, nil
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
					a.transferBlockFiles(chainID)
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

func (a *ArchiveService) handleMessages() {
	logger.Infof("Archive service - handle gossip msg for peer %s: Entering", a.peerId)
	defer logger.Infof("Archive service - handle gossip msg for peer %s: Exiting", a.peerId)
	msgChan := a.gossip.Accept()
	for {
		select {
		case <-a.stopCh:
			return
		case msg := <-msgChan:
			a.handleMessage(msg)
		}
	}
}

func (a *ArchiveService) handleMessage(msg *proto.GossipMessage) {
	a.lock.Lock()
	defer a.lock.Unlock()

	chainID := string(msg.Channel)
	ledger, err := a.getLedger(chainID)
	if err != nil {
		logger.Errorf("Archive service - get ledger for chain: %s failed with error: %s", chainID, err)
		return
	}

	metaInfo := msg.GetArchiveMsg().ArchiveMetaInfo
	logger.Infof("Archive service - updating meta info[%s] for channel: %s", spew.Sdump(metaInfo), chainID)
	ledger.UpdateArchiveMetaInfo(metaInfo)
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

func (a *ArchiveService) Stop() {
	a.gossip.Stop()
	close(a.stopCh)

	for _, w := range a.watchers {
		if err := w.Close(); err != nil {
			logger.Warnf("Stop watcher[%s] got error: %s", spew.Sdump(w), err)
			continue
		}
	}

	if err := a.dfsClient.Close(); err != nil {
		logger.Warnf("Stop dfs client[%s] got error: %s", spew.Sdump(a.dfsClient), err)
	}
}

func getLedgerDir(chainID string) string {
	rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
	chainsDir := filepath.Join(kvledger.BlockStorePath(rootFSPath), hybridblkstorage.ChainsDir)
	ledgerDir := filepath.Join(chainsDir, chainID)
	return ledgerDir
}
