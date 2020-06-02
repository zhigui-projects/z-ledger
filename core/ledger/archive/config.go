package archive

import (
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
)

// Config is a structure used to configure the archive service for the ledger.
type Config struct {
	Enabled  bool
	Type     string
	HdfsConf *common.HdfsConfig
	IpfsConf *common.IpfsConfig
}
