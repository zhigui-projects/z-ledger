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
	hs "github.com/zhigui-projects/go-hotstuff/consensus"
	"github.com/zhigui-projects/go-hotstuff/pacemaker"
	"github.com/zhigui-projects/go-hotstuff/pb"
)

var (
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
)

type chain struct {
	hsb      *hs.HotStuffBase
	pm       pacemaker.PaceMaker
	sendChan chan *message
	exitChan chan struct{}
	cancel   context.CancelFunc
	support  consensus.ConsenterSupport
	logger   *flogging.FabricLogger
	SubmitClient
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
		seq := c.support.Sequence()
		err = nil
		select {
		case msg := <-c.sendChan:
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
					go c.submit([][]byte{cmds})
				}
				if len(batches) > 0 {
					pmTimer = time.After(100 * time.Millisecond)
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
					go c.submit([][]byte{cmds})

					pmTimer = time.After(200 * time.Millisecond)
				}

				cmdsReq := &CmdsRequest{
					MsgType: ConfigMsg,
					Batche:  []*cb.Envelope{msg.configMsg},
				}
				cmds, _ := json.Marshal(cmdsReq)
				go func() {
					c.pm.Submit(cmds)
					c.pm.Submit(nil)
					c.pm.Submit(nil)
					c.pm.Submit(nil)
				}()

				timer = nil
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
			go c.submit([][]byte{cmds})

			pmTimer = time.After(200 * time.Millisecond)
		case <-pmTimer:
			go func() {
				c.submit([][]byte{nil, nil, nil})
			}()
		case <-c.exitChan:
			c.logger.Debugf("Exiting")
			return
		}
	}
}

func (c *chain) submit(txs [][]byte) {
	chs := c.hsb.GetHotStuff(c.support.ChannelID())
	leader := chs.GetLeader(chs.GetCurView())
	hsc := c.getSubmitClient(leader)
	for _, cmds := range txs {
		if hsc == nil {
			c.pm.Submit(cmds)
		} else {
			_, err := hsc.Submit(context.Background(), &pb.SubmitRequest{Chain: c.support.ChannelID(), Cmds: cmds})
			if err != nil {
				c.pm.Submit(cmds)
				c.delSubmitClient(leader)
			}
		}
	}
}
