package hotstuff

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	hs "github.com/hyperledger/fabric-protos-go/orderer/hotstuff"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/orderer/common/localconfig"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/pkg/errors"
	"github.com/zhigui-projects/go-hotstuff/common/crypto"
	hsc "github.com/zhigui-projects/go-hotstuff/consensus"
	"github.com/zhigui-projects/go-hotstuff/pacemaker"
	"github.com/zhigui-projects/go-hotstuff/pb"
)

var logger = flogging.MustGetLogger("orderer.consensus.hotstuff")

type consenter struct {
	nodeId        int64
	cert          []byte
	hsb           *hsc.HotStuffBase
	md            *pb.ConfigMetadata
	ordererConfig localconfig.TopLevel
}

// New creates a hotstuff Consenter
func New(conf *localconfig.TopLevel, srvConf comm.ServerConfig) consensus.Consenter {
	return &consenter{
		cert:          srvConf.SecOpts.Certificate,
		ordererConfig: *conf,
	}
}

func (c *consenter) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	logger.Infof("Starting a chain: %s", support.ChannelID())

	if c.hsb == nil {
		hsb, err := c.newHotStuff(support)
		if err != nil {
			return nil, err
		}
		c.hsb = hsb
	}

	decideExec := func(cmds []byte) {
		if len(cmds) == 0 {
			return
		}
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

	return &chain{
		hsb:           c.hsb,
		pm:            pacemaker.NewRoundRobinPM(c.hsb.GetHotStuff(support.ChannelID()), c.nodeId, c.md, decideExec),
		sendChan:      make(chan *message),
		support:       support,
		submitClients: make(map[int64]pb.HotstuffClient),
	}, nil
}

func (c *consenter) newHotStuff(support consensus.ConsenterSupport) (*hsc.HotStuffBase, error) {
	var selfId int64 = -1
	_, selfDer, err := validateCert(c.cert)
	if err != nil {
		return nil, err
	}

	m := &hs.ConfigMetadata{}
	if err := proto.Unmarshal(support.SharedConfig().HotStuffMetadata(), m); err != nil {
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
	c.md = replicas.Metadata

	nodes := make([]*hsc.NodeInfo, len(m.Consenters))
	for i, consenter := range m.Consenters {
		// TODO tls
		node := &hsc.NodeInfo{
			Id:   hsc.ReplicaID(i),
			Addr: consenter.Host + ":" + strconv.FormatUint(uint64(consenter.Port), 10),
		}
		nodes[int64(i)] = node

		pk, _, err := validateCert(consenter.SignCert)
		if err != nil {
			logger.Errorf("Failed converting PEM to public key, err: %v, cert: %s", err, string(consenter.SignCert))
			return nil, err
		}
		pubkey, ok := pk.(*ecdsa.PublicKey)
		if !ok {
			logger.Error("Invalid key type. It should be *ecdsa.PublicKey")
			return nil, errors.New("Invalid key type. It should be *ecdsa.PublicKey")
		}

		replicas.Replicas[hsc.ReplicaID(i)] = &hsc.ReplicaInfo{Verifier: &crypto.ECDSAVerifier{Pub: pubkey}}

		if _, der, err := validateCert(consenter.ServerTlsCert); err != nil {
			return nil, err
		} else if bytes.Equal(selfDer, der) {
			selfId = int64(i)
		}
	}

	if selfId == -1 {
		return nil, errors.Errorf("Could not find cert:%s among consenters", string(c.cert))
	}
	c.nodeId = selfId

	keyCert, err := loadX509KeyPair(c.ordererConfig.General.LocalMSPDir)
	if err != nil {
		return nil, err
	}
	priKey := keyCert.PrivateKey.(*ecdsa.PrivateKey)

	return hsc.NewHotStuffBase(hsc.ReplicaID(selfId), nodes, &crypto.ECDSASigner{Pri: priKey}, replicas), nil
}

func validateCert(pemData []byte) (interface{}, []byte, error) {
	bl, _ := pem.Decode(pemData)
	if bl == nil {
		return nil, nil, errors.Errorf("Certificate is not PEM encoded: %s", string(pemData))
	}

	if cert, err := x509.ParseCertificate(bl.Bytes); err != nil {
		return nil, nil, errors.Errorf("Certificate has invalid ASN1 structure, err:%v, data:%s", err, string(pemData))
	} else {
		return cert.PublicKey, bl.Bytes, nil
	}
}

func loadX509KeyPair(dir string) (*tls.Certificate, error) {
	signPath := filepath.Join(dir, "signcerts")
	fis, err := ioutil.ReadDir(signPath)
	if err != nil {
		return nil, errors.Errorf("ReadDir signcerts failed, err: %v", err)
	}
	if len(fis) != 1 {
		return nil, errors.Errorf("Invalid signcerts dir there are %v files", len(fis))
	}
	signFile := filepath.Join(signPath, fis[0].Name())

	keyPath := filepath.Join(dir, "keystore")
	keyFis, err := ioutil.ReadDir(keyPath)
	if err != nil {
		return nil, errors.Errorf("ReadDir keystore failed, err: %v", err)
	}
	if len(keyFis) != 1 {
		return nil, errors.Errorf("Invalid keystore dir there are %v files", len(keyFis))
	}
	keyFile := filepath.Join(keyPath, keyFis[0].Name())

	cert, err := tls.LoadX509KeyPair(signFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
