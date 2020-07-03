package consensus

import (
	"context"
	"encoding/hex"
	"sync"

	"github.com/pkg/errors"
	"github.com/zhigui-projects/go-hotstuff/common/log"
	"github.com/zhigui-projects/go-hotstuff/pacemaker"
	"github.com/zhigui-projects/go-hotstuff/pb"
	"google.golang.org/grpc/peer"
)

type ChainHotStuff struct {
	pacemaker.PaceMaker
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

func (c *ChainHotStuff) ApplyPaceMaker(pm pacemaker.PaceMaker) {
	c.PaceMaker = pm
}

func (c *ChainHotStuff) DoVote(leader int64, vote *pb.Vote) {
	if leader != c.HotStuffCore.GetID() {
		if err := c.UnicastMsg(&pb.Message{Chain: c.chainId, Type: &pb.Message_Vote{Vote: vote}}, leader); err != nil {
			c.logger.Error("do vote error when unicast msg", "to", leader)
		}
	} else {
		if err := c.HotStuffCore.OnReceiveVote(vote); err != nil {
			c.logger.Warning("do vote error when receive vote", "to", leader)
		}
	}
}

func (c *ChainHotStuff) GetChainId() string {
	return c.chainId
}

type HotStuffBase struct {
	hscs map[string]*ChainHotStuff
	*NodeManager
	queue chan MsgExecutor

	rwmutex  sync.RWMutex
	id       ReplicaID
	signer   Signer
	replicas *ReplicaConf
	logger   log.Logger
}

func NewHotStuffBase(id ReplicaID, nodes []*NodeInfo, signer Signer, replicas *ReplicaConf) *HotStuffBase {
	logger := log.GetLogger("module", "consensus")
	if len(nodes) == 0 {
		logger.Error("not found hotstuff replica node info")
		return nil
	}

	hsb := &HotStuffBase{
		hscs:        make(map[string]*ChainHotStuff),
		NodeManager: NewNodeManager(id, nodes, logger),
		queue:       make(chan MsgExecutor),
		id:          id,
		signer:      signer,
		replicas:    replicas,
		logger:      logger,
	}
	pb.RegisterHotstuffServer(hsb.Server(), hsb)
	return hsb
}

func (hsb *HotStuffBase) ApplyPaceMaker(pm pacemaker.PaceMaker, chain string) {
	hsb.GetHotStuff(chain).ApplyPaceMaker(pm)
}

func (hsb *HotStuffBase) handleProposal(proposal *pb.Proposal, chain string) {
	if proposal == nil || proposal.Block == nil {
		hsb.logger.Warning("handle proposal with empty block")
		return
	}
	hsb.logger.Info("handle proposal", "proposer", proposal.Block.Proposer,
		"height", proposal.Block.Height, "hash", hex.EncodeToString(proposal.Block.SelfQc.BlockHash))

	if err := hsb.GetHotStuff(chain).HotStuffCore.OnReceiveProposal(proposal); err != nil {
		hsb.logger.Warning("handle proposal catch error", "error", err)
	}
}

func (hsb *HotStuffBase) handleVote(vote *pb.Vote, chain string) {
	if ok := hsb.replicas.VerifyVote(vote); !ok {
		return
	}
	if err := hsb.GetHotStuff(chain).OnReceiveVote(vote); err != nil {
		hsb.logger.Warning("handle vote catch error", "error", err)
		return
	}
}

func (hsb *HotStuffBase) handleNewView(id ReplicaID, newView *pb.NewView, chain string) {
	block, err := hsb.GetHotStuff(chain).getBlockByHash(newView.GenericQc.BlockHash)
	if err != nil {
		hsb.logger.Error("Could not find block of new QC", "error", err)
		return
	}

	//hsb.UpdateHighestQC(block, newView.GenericQc)

	hsb.hscs[chain].notify(&ReceiveNewViewEvent{int64(id), block, newView})
}

func (hsb *HotStuffBase) DoVote(leader int64, vote *pb.Vote, chain string) {
	hsb.GetHotStuff(chain).DoVote(leader, vote)
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

	if err := hsb.hscs[req.Chain].Submit(req.Cmds); err != nil {
		return &pb.SubmitResponse{Status: pb.Status_SERVICE_UNAVAILABLE}, err
	} else {
		return &pb.SubmitResponse{Status: pb.Status_SUCCESS}, nil
	}
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
	hsc := NewHotStuffCore(hsb.id, hsb.signer, hsb.replicas, hsb.logger)
	chs = &ChainHotStuff{HotStuffCore: hsc, NodeManager: hsb.NodeManager, chainId: chain}
	hsb.hscs[chain] = chs
	return chs
}

func (hsb *HotStuffBase) receiveMsg(msg *pb.Message, src ReplicaID) {
	hsb.logger.Debug("received message", "from", src, "to", hsb.hscs[msg.Chain].GetID(), "msg", msg.String())

	switch msg.GetType().(type) {
	case *pb.Message_Proposal:
		hsb.handleProposal(msg.GetProposal(), msg.Chain)
	case *pb.Message_NewView:
		hsb.handleNewView(src, msg.GetNewView(), msg.Chain)
	case *pb.Message_Vote:
		hsb.handleVote(msg.GetVote(), msg.Chain)
	}
}
