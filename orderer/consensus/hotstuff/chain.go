/*
Copyright Zhigui.com Corp. 2020 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hotstuff

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/pkg/errors"
	"github.com/zhigui-projects/go-hotstuff/api"
	hs "github.com/zhigui-projects/go-hotstuff/consensus"
	"github.com/zhigui-projects/go-hotstuff/pb"
)

var (
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
)

type chain struct {
	hsb        *hs.HotStuffBase
	pm         api.PaceMaker
	sendChan   chan *message
	exitChan   chan struct{}
	cancel     context.CancelFunc
	support    consensus.ConsenterSupport
	logger     *flogging.FabricLogger
	applyC     chan []byte
	pmWaitNsec int64
}

type message struct {
	configSeq uint64
	normalMsg *cb.Envelope
	configMsg *cb.Envelope
}

const (
	NormalMsg int8 = 1
	ConfigMsg      = 2
)

type CmdsRequest struct {
	MsgType int8
	Batche  []*cb.Envelope
}

// Start instructs the orderer to begin serving the chain and keep it current.
func (c *chain) Start() {
	once.Do(func() {
		ctx, cancel = context.WithCancel(context.Background())
		go c.hsb.Start(ctx)
	})
	c.cancel = cancel
	go c.pm.Run(ctx)
	go c.main()
}

// Halt stops the chain.
func (c *chain) Halt() {
	select {
	case <-c.exitChan:
		// Allow multiple halts without panic
	default:
		close(c.exitChan)
	}
	// stop hotstuff
	c.cancel()
}

// Order submits normal type transactions for ordering.
func (c *chain) Order(env *common.Envelope, configSeq uint64) error {
	// TODO: Delete when benchmark
	chr, _ := unmarshalEnvelope(env)
	c.logger.Infof("Receive normal transaction, txId: %s", chr.TxId)

	select {
	case c.sendChan <- &message{
		configSeq: configSeq,
		normalMsg: env,
	}:
		return nil
	case <-c.exitChan:
		return errors.Errorf("Exiting")
	}
}

// Configure submits config type transactions for ordering.
func (c *chain) Configure(config *common.Envelope, configSeq uint64) error {
	// TODO: Delete when benchmark
	chr, _ := unmarshalEnvelope(config)
	c.logger.Infof("Receive config transaction, txId: %s", chr.TxId)

	select {
	case c.sendChan <- &message{
		configSeq: configSeq,
		configMsg: config,
	}:
		return nil
	case <-c.exitChan:
		return errors.Errorf("Exiting")
	}
}

func (c *chain) WaitReady() error {
	return nil
}

// Errored returns a channel that closes when the chain stops.
func (c *chain) Errored() <-chan struct{} {
	return c.exitChan
}

func (c *chain) main() {
	var timer <-chan time.Time
	var pmTimer <-chan time.Time
	var err error

	for {
		select {
		case msg := <-c.sendChan:
			seq := c.support.Sequence()
			if msg.configMsg == nil {
				// NormalMsg
				if msg.configSeq < seq {
					_, err = c.support.ProcessNormalMsg(msg.normalMsg)
					if err != nil {
						c.logger.Warningf("Discarding bad normal message: %s", err)
						continue
					}
				}
				batches, pending := c.support.BlockCutter().Ordered(msg.normalMsg)

				for _, batch := range batches {
					cmdsReq := &CmdsRequest{
						MsgType: NormalMsg,
						Batche:  batch,
					}
					cmds, _ := json.Marshal(cmdsReq)
					go c.pm.Submit(cmds)
				}
				if len(batches) > 0 {
					pmTimer = time.After(time.Duration(c.pmWaitNsec))
				}

				switch {
				case timer != nil && !pending:
					// Timer is already running but there are no messages pending, stop the timer
					timer = nil
				case timer == nil && pending:
					// Timer is not already running and there are messages pending, so start it
					timer = time.After(c.support.SharedConfig().BatchTimeout())
					c.logger.Debugf("Just began %s batch timer", c.support.SharedConfig().BatchTimeout().String())
				default:
					// Do nothing when:
					// 1. Timer is already running and there are messages pending
					// 2. Timer is not set and there are no messages pending
				}

			} else {
				// ConfigMsg
				if msg.configSeq < seq {
					msg.configMsg, _, err = c.support.ProcessConfigMsg(msg.configMsg)
					if err != nil {
						c.logger.Warningf("Discarding bad config message: %s", err)
						continue
					}
				}
				batch := c.support.BlockCutter().Cut()
				if batch != nil {
					cmdsReq := &CmdsRequest{
						MsgType: NormalMsg,
						Batche:  batch,
					}
					cmds, _ := json.Marshal(cmdsReq)
					go c.pm.Submit(cmds)

					pmTimer = time.After(time.Duration(c.pmWaitNsec))
				}

				cmdsReq := &CmdsRequest{
					MsgType: ConfigMsg,
					Batche:  []*cb.Envelope{msg.configMsg},
				}
				cmds, _ := json.Marshal(cmdsReq)
				go func() {
					// submit more nil tx to avoid real tx not exec, keep alive
					c.pm.Submit(cmds)
					c.pm.Submit(nil)
					c.pm.Submit(nil)
					c.pm.Submit(nil)
				}()
				go c.submit(nil)

				timer = nil
			}
		case cmds := <-c.applyC:
			var req CmdsRequest
			if err := json.Unmarshal(cmds, &req); err != nil {
				c.logger.Errorf("cmds request unmarshal failed, err:%v", err)
				return
			}

			// TODO: Delete when benchmark
			//for _, env := range req.Batche {
			//	chr, _ := unmarshalEnvelope(env)
			//	c.logger.Infof("Consensus complete for transaction, txId: %s", chr.TxId)
			//}

			block := c.support.CreateNextBlock(req.Batche)
			if req.MsgType == NormalMsg {
				c.support.WriteBlock(block, nil)
				c.logger.Infof("Writing block [%d] to ledger", block.Header.Number)
			} else {
				c.support.WriteConfigBlock(block, nil)
				c.logger.Infof("Writing config block [%d] to ledger", block.Header.Number)
			}
		case <-timer:
			//clear the timer
			timer = nil

			batch := c.support.BlockCutter().Cut()
			if len(batch) == 0 {
				c.logger.Warningf("Batch timer expired with no pending requests, this might indicate a bug")
				continue
			}
			c.logger.Debugf("Batch timer expired, creating block")

			cmdsReq := &CmdsRequest{
				MsgType: NormalMsg,
				Batche:  batch,
			}
			cmds, _ := json.Marshal(cmdsReq)
			go c.pm.Submit(cmds)

			pmTimer = time.After(time.Duration(c.pmWaitNsec))
		case <-pmTimer:
			go func() {
				c.pm.Submit(nil)
				c.pm.Submit(nil)
				c.pm.Submit(nil)
			}()

			go c.submit(nil)
		case <-c.exitChan:
			c.logger.Debugf("Exiting")
			return
		}
	}
}

func (c *chain) submit(cmds []byte) {
	chs := c.hsb.GetHotStuff(c.support.ChannelID())
	leader := chs.GetLeader(chs.GetCurView())
	if leader != c.hsb.GetID() {
		go c.hsb.UnicastMsg(&pb.Message{Chain: c.support.ChannelID(), Type: &pb.Message_Forward{
			Forward: &pb.Forward{Data: cmds}}}, leader)
	}
}
