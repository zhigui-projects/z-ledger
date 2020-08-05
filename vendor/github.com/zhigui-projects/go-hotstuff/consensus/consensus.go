/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package consensus

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/zhigui-projects/go-hotstuff/api"
	"github.com/zhigui-projects/go-hotstuff/common/log"
	"github.com/zhigui-projects/go-hotstuff/pb"
	"google.golang.org/grpc/peer"
)

type ChainHotStuff struct {
	api.PaceMaker
	*HotStuffCore
	*NodeManager
	chainId string
}

func (c *ChainHotStuff) Start(ctx context.Context) {
	for {
		select {
		case n := <-c.HotStuffCore.GetNotifier():
			go n.ExecuteEvent(c.PaceMaker)
		case <-ctx.Done():
			return
		}
	}
}

func (c *ChainHotStuff) ApplyPaceMaker(pm api.PaceMaker) {
	c.PaceMaker = pm
}

func (c *ChainHotStuff) DoVote(leader int64, vote *pb.Vote) {
	if leader != c.HotStuffCore.GetID() {
		if err := c.UnicastMsg(&pb.Message{Chain: c.chainId, Type: &pb.Message_Vote{Vote: vote}}, leader); err != nil {
			c.logger.Error("do vote error when unicast msg", "to", leader)
		}
	} else {
		if err := c.HotStuffCore.OnProposalVote(vote); err != nil {
			c.logger.Warning("do vote error when receive vote", "to", leader)
		}
	}
}

func (c *ChainHotStuff) GetChainId() string {
	return c.chainId
}

type HotStuffBase struct {
	api.PaceMaker
	*HotStuffCore
	*NodeManager
	queue  chan MsgExecutor
	logger api.Logger

	hscs    map[string]*ChainHotStuff
	rwmutex sync.RWMutex
}

func NewHotStuffBase(id ReplicaID, nodes []*NodeInfo, signer api.Signer, replicas *ReplicaConf) *HotStuffBase {
	logger := log.GetLogger("module", "consensus", "node", id)
	if len(nodes) == 0 {
		logger.Error("not found hotstuff replica node info")
		return nil
	}

	hsb := &HotStuffBase{
		HotStuffCore: NewHotStuffCore(id, signer, replicas, logger),
		NodeManager:  NewNodeManager(id, nodes, logger),
		queue:        make(chan MsgExecutor),
		logger:       logger,
		hscs:         make(map[string]*ChainHotStuff),
	}
	pb.RegisterHotstuffServer(hsb.Server(), hsb)
	return hsb
}

func (hsb *HotStuffBase) GetHotStuff(chain string) *ChainHotStuff {
	hsb.rwmutex.RLock()
	chs := hsb.hscs[chain]
	hsb.rwmutex.RUnlock()
	if chs != nil {
		return chs
	}

	hsb.rwmutex.Lock()
	defer hsb.rwmutex.Unlock()
	hsc := NewHotStuffCore(hsb.id, hsb.Signer, hsb.replicas, hsb.logger)
	chs = &ChainHotStuff{HotStuffCore: hsc, NodeManager: hsb.NodeManager, chainId: chain}
	hsb.hscs[chain] = chs
	return chs
}

func (hsb *HotStuffBase) ApplyPaceMaker(pm api.PaceMaker, chain string) {
	hsb.GetHotStuff(chain).ApplyPaceMaker(pm)
}

func (hsb *HotStuffBase) Start(ctx context.Context) {
	go hsb.StartServer()
	go hsb.ConnectWorkers(hsb.queue)

	for {
		select {
		case m := <-hsb.queue:
			go m.ExecuteMessage(hsb)
		case <-ctx.Done():
			return
		}
	}
}

func (hsb *HotStuffBase) Submit(ctx context.Context, req *pb.SubmitRequest) (*pb.SubmitResponse, error) {
	var remoteAddress string
	if p, _ := peer.FromContext(ctx); p != nil {
		if address := p.Addr; address != nil {
			remoteAddress = address.String()
		}
	}
	hsb.logger.Info("receive new submit request", "remoteAddress", remoteAddress)
	if len(req.Cmds) == 0 {
		return &pb.SubmitResponse{Status: pb.Status_BAD_REQUEST}, errors.New("request data is empty")
	}
	if req.Chain == "" {
		return &pb.SubmitResponse{Status: pb.Status_BAD_REQUEST}, errors.New("chainId is empty")
	}

	if err := hsb.GetHotStuff(req.Chain).Submit(req.Cmds); err != nil {
		return &pb.SubmitResponse{Status: pb.Status_SERVICE_UNAVAILABLE}, err
	} else {
		return &pb.SubmitResponse{Status: pb.Status_SUCCESS}, nil
	}
}

func (hsb *HotStuffBase) handleMessage(src ReplicaID, msg *pb.Message) {
	hsb.logger.Debug("received message", "from", src, "msg", msg.String())
	if msg.Chain == "" {
		hsb.logger.Error("chainId is empty")
		return
	}

	switch msg.GetType().(type) {
	case *pb.Message_Proposal:
		hsb.handleProposal(msg.GetProposal(), msg.Chain)
	case *pb.Message_ViewChange:
		hsb.handleViewChange(src, msg.GetViewChange(), msg.Chain)
	case *pb.Message_Vote:
		hsb.handleVote(src, msg.GetVote(), msg.Chain)
	case *pb.Message_Forward:
		hsb.handleForward(src, msg.GetForward(), msg.Chain)
	default:
		hsb.logger.Error("received invalid message type", "from", src, "msg", msg.String())
	}
}

// TODO 封装on接口
func (hsb *HotStuffBase) handleProposal(proposal *pb.Proposal, chain string) {
	if proposal.Block == nil {
		hsb.logger.Error("receive invalid proposal msg with empty block")
		return
	}

	if err := hsb.GetHotStuff(chain).HotStuffCore.OnReceiveProposal(proposal); err != nil {
		hsb.logger.Error("handle proposal msg catch error", "error", err)
	}
}

func (hsb *HotStuffBase) handleViewChange(id ReplicaID, vc *pb.ViewChange, chain string) {
	if ok := hsb.GetHotStuff(chain).GetReplicas().VerifyIdentity(id, vc.Signature, vc.Data); !ok {
		hsb.logger.Error("receive invalid view change msg signature verify failed")
		return
	}
	newView := &pb.NewView{}
	if err := proto.Unmarshal(vc.Data, newView); err != nil {
		hsb.logger.Error("receive invalid view change msg that unmarshal NewView failed", "from", id, "error", err)
		return
	}
	if newView.GenericQc == nil {
		hsb.logger.Error("receive invalid view change msg that genericQc is null", "from", id)
		return
	}

	if err := hsb.GetHotStuff(chain).OnNewView(int64(id), newView); err != nil {
		hsb.logger.Error("handle view change msg catch error", "from", id, "error", err)
		return
	}
}

func (hsb *HotStuffBase) handleVote(id ReplicaID, vote *pb.Vote, chain string) {
	if len(vote.BlockHash) == 0 || vote.Cert == nil {
		hsb.logger.Error("receive invalid vote msg data", "from", id)
		return
	}
	if ok := hsb.GetHotStuff(chain).GetReplicas().VerifyIdentity(id, vote.Cert.Signature, vote.BlockHash); !ok {
		hsb.logger.Error("receive invalid vote msg signature verify failed")
		return
	}
	if int64(id) != vote.Voter {
		hsb.logger.Error("receive invalid vote msg replica id not match", "from", id, "voter", vote.Voter)
		return
	}
	if leader := hsb.GetHotStuff(chain).GetLeader(vote.ViewNumber); ReplicaID(leader) != hsb.GetHotStuff(chain).id {
		hsb.logger.Error("receive invalid vote msg not match leader id", "from", id,
			"view", vote.ViewNumber, "leader", leader)
		return
	}

	if err := hsb.GetHotStuff(chain).OnProposalVote(vote); err != nil {
		hsb.logger.Error("handle vote msg catch error", "from", id, "error", err)
		return
	}
}

func (hsb *HotStuffBase) handleForward(id ReplicaID, msg *pb.Forward, chain string) {
	if len(msg.Data) == 0 {
		hsb.logger.Error("receive invalid forward msg data", "from", id)
		return
	}
	_ = hsb.GetHotStuff(chain).PaceMaker.Submit(msg.Data)
}
