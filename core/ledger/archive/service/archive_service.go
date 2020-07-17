/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	proto "github.com/hyperledger/fabric-protos-go/gossip"
	"github.com/hyperledger/fabric/common/flogging"
	le "github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/archive"
	"github.com/hyperledger/fabric/core/ledger/dfs"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/hyperledger/fabric/core/ledger/ledgermgmt"
	"github.com/pkg/errors"
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
	// StartTickerForChannel dynamically starts ticker of ledger files.
	StartTickerForChannel(chainID string) error

	// StopTickerForChannel dynamically stops ticker of ledger files.
	StopTickerForChannel(chainID string) error

	// Stop stops ticker of ledger files.
	Stop()
}

type peerID []byte

type ArchiveService struct {
	archiveSvc
	gossip    ArchiveGossip
	id        peerID
	ledgerMgr *ledgermgmt.LedgerMgr
	dfsClient common.FsClient
	tickers   map[string]*time.Ticker
	lock      sync.RWMutex
	stopCh    chan struct{}
	conf      *archive.Trigger
}

// New construction function to create and initialize
// archive service instance. It tries to establish connection to
// the specified dfs name node, in case it fails to dial to it, return nil
func New(g gossip, ledgerMgr *ledgermgmt.LedgerMgr, channel string, peerId string, conf *archive.Config) (*ArchiveService, error) {
	dfsConf := &common.Config{Type: conf.Type, HdfsConf: conf.HdfsConf, IpfsConf: conf.IpfsConf}
	client, err := dfs.NewDfsClient(dfsConf)
	if err != nil {
		logger.Errorf("Archive service can't connect to HDFS, due to %+v", err)
		return nil, err
	}

	a := &ArchiveService{
		gossip:    NewGossipImpl(g, channel),
		id:        peerID(peerId),
		ledgerMgr: ledgerMgr,
		dfsClient: client,
		tickers:   make(map[string]*time.Ticker),
		stopCh:    make(chan struct{}),
		conf:      conf.Trigger,
	}
	go a.handleMessages()

	return a, nil
}

// StartTickerForChannel starts a ticker to trigger the block files transfer for channel
func (a *ArchiveService) StartTickerForChannel(chainID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, exist := a.tickers[chainID]; exist {
		errMsg := fmt.Sprintf("Archive service - ledger ticker already exists for %s found, can't start a new ticker", chainID)
		logger.Warn(errMsg)
		return errors.New(errMsg)
	}
	ticker := time.NewTicker(a.conf.CheckInterval)

	go a.broadcastMetaInfo(chainID)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				logger.Debugf("Archive service - ticker tick at: %s", t)
				if t.After(a.conf.BeginTime) && t.Before(a.conf.EndTime) {
					logger.Infof("Archive service - block files transfer triggered at: %s", t)
					a.transferBlockFiles(chainID)
					go a.broadcastMetaInfo(chainID)
				}
			}
		}
	}()

	a.tickers[chainID] = ticker
	return nil
}

func (a *ArchiveService) StopTickerForChannel(chainID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	ticker, exist := a.tickers[chainID]
	if !exist {
		errMsg := fmt.Sprintf("Archive service - ledger ticker does not exists for %s found", chainID)
		logger.Warn(errMsg)
		return errors.New(errMsg)
	}

	ticker.Stop()
	return nil
}

func (a *ArchiveService) Stop() {
	a.gossip.Stop()
	close(a.stopCh)

	for _, t := range a.tickers {
		t.Stop()
	}

	if err := a.dfsClient.Close(); err != nil {
		logger.Warnf("Stop dfs client[%s] got error: %s", spew.Sdump(a.dfsClient), err)
	}
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
	logger.Infof("Archive service - handle gossip msg for peer[%s]: Entering", string(a.id))
	defer logger.Infof("Archive service - handle gossip msg for peer[%s]: Exiting", string(a.id))
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
