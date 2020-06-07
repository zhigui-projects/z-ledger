package dfs

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/hyperledger/fabric/core/ledger/dfs/hdfs"
	"github.com/hyperledger/fabric/core/ledger/dfs/ipfs"
)

const (
	HDFS = "hdfs"
	IPFS = "ipfs"
)

var logger = flogging.MustGetLogger("dfs")

func NewDfsClient(conf *common.Config) (common.FsClient, error) {
	logger.Infof("initializing dfs client with config: %s", spew.Sdump(conf))

	var fsClient common.FsClient
	var err error
	if conf == nil {
		return nil, errors.New("can't new dfs client due to archive config is nil")
	}

	if conf.Type == HDFS {
		if fsClient, err = hdfs.NewFsClient(conf.HdfsConf); err != nil {
			return nil, err
		}
	} else if conf.Type == IPFS {
		if fsClient, err = ipfs.NewFsClient(conf.IpfsConf); err != nil {
			return nil, err
		}
	} else {
		errMsg := fmt.Sprintf("Unknown type of dfs, expected [hdfs or ipfs], got [%s]", conf.Type)
		logger.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	return fsClient, nil
}
