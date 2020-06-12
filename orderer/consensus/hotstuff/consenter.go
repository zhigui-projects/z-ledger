package hotstuff

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strconv"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	hs "github.com/hyperledger/fabric-protos-go/orderer/hotstuff"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/pkg/errors"
	"github.com/zhigui-projects/go-hotstuff/common/crypto"
	hsc "github.com/zhigui-projects/go-hotstuff/consensus"
	"github.com/zhigui-projects/go-hotstuff/pacemaker"
	"github.com/zhigui-projects/go-hotstuff/pb"
)

var logger = flogging.MustGetLogger("orderer.consensus.hotstuff")

type consenter struct {
	cert []byte
}

// New creates a hotstuff Consenter
func New(srvConf comm.ServerConfig) consensus.Consenter {
	return &consenter{
		cert: srvConf.SecOpts.Certificate,
	}
}

func (c *consenter) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	logger.Infof("Starting a chain: %s", support.ChannelID())

	var selfId int64 = -1
	selfDer, err := validateCert(c.cert)
	if err != nil {
		return nil, err
	}

	m := &hs.ConfigMetadata{}
	if err := proto.Unmarshal(support.SharedConfig().ConsensusMetadata(), m); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal consensus metadata")
	}
	if m.Options == nil {
		return nil, errors.New("hotstuff options have not been provided")
	}
	if len(m.Consenters) == 0 {
		return nil, errors.New("hotstuff consenters have not been provided")
	}

	replicas := &hsc.ReplicaConf{
		Metadata: &pb.ConfigMetadata{
			N:              m.Options.N,
			F:              m.Options.F,
			MsgWaitTimeout: m.Options.RequestTimeoutSec,
		},
		Replicas: make(map[hsc.ReplicaID]*hsc.ReplicaInfo, len(m.Consenters)),
	}
	nodes := make([]*hsc.NodeInfo, len(m.Consenters))
	for i, consenter := range m.Consenters {
		// TODO tls
		node := &hsc.NodeInfo{
			Id:   hsc.ReplicaID(i),
			Addr: consenter.Host + ":" + strconv.FormatUint(uint64(consenter.Port), 10),
		}
		nodes[int64(i)] = node

		pk, err := utils.PEMtoPublicKey(consenter.SignCert, nil)
		if err != nil {
			logger.Errorf("Failed converting PEM to public key, err: %v", err)
			return nil, err
		}
		pubkey, ok := pk.(*ecdsa.PublicKey)
		if !ok {
			logger.Error("Invalid key type. It should be *ecdsa.PublicKey")
			return nil, errors.New("Invalid key type. It should be *ecdsa.PublicKey")
		}

		replicas.Replicas[hsc.ReplicaID(i)] = &hsc.ReplicaInfo{Verifier: &crypto.ECDSAVerifier{Pub: pubkey}}

		if der, err := validateCert(consenter.ServerTlsCert); err != nil {
			return nil, err
		} else if bytes.Equal(selfDer, der) {
			selfId = int64(i)
		}
	}

	if selfId == -1 {
		return nil, errors.Errorf("Could not find cert:%s among consenters", string(c.cert))
	}

	decideExec := func(cmds []byte) {
		var req CmdsRequest
		if err := json.Unmarshal(cmds, &req); err != nil {
			logger.Errorf("cmds request unmarshal failed, err:%v", err)
			return
		}

		block := support.CreateNextBlock(req.Batche)
		if req.MsgType == NormalMsg {
			support.WriteBlock(block, nil)
			logger.Infof("Writing block [%d] for channelId: %s to ledger", block.Header.Number, support.ChannelID())
		} else {
			support.WriteConfigBlock(block, nil)
			logger.Infof("Writing config block [%d] for channelId: %s to ledger", block.Header.Number, support.ChannelID())
		}
	}

	hsb := hsc.NewHotStuffBase(hsc.ReplicaID(selfId), nodes, support, replicas)
	return &chain{
		pm:       pacemaker.NewRoundRobinPM(hsb, selfId, replicas.Metadata, decideExec),
		sendChan: make(chan *message),
		support:  support,
	}, nil
}

func validateCert(pemData []byte) ([]byte, error) {
	bl, _ := pem.Decode(pemData)
	if bl == nil {
		return nil, errors.Errorf("Certificate is not PEM encoded: %s", string(pemData))
	}

	if _, err := x509.ParseCertificate(bl.Bytes); err != nil {
		return nil, errors.Errorf("Certificate has invalid ASN1 structure, err:%v, data:%s", err, string(pemData))
	}
	return bl.Bytes, nil
}
