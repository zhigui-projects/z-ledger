package archive

import (
	"time"

	"github.com/hyperledger/fabric/core/ledger/dfs/common"
)

type Trigger struct {
	CheckInterval time.Duration
	BeginTime     time.Time
	EndTime       time.Time
}

// Config is a structure used to configure the archive service for the ledger.
type Config struct {
	Enabled  bool
	Trigger  *Trigger
	Type     string
	HdfsConf *common.HdfsConfig
	IpfsConf *common.IpfsConfig
	FsRoot   string
}
