package ipfs

import (
	"os"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
)

var logger = flogging.MustGetLogger("dfs.ipfs")

type FsClient struct {
}

func NewFsClient(conf *common.IpfsConfig) (*FsClient, error) {
	return nil, nil
}

func (c *FsClient) ReadDir(dirname string) ([]os.FileInfo, error) {
	return nil, nil
}

func (c *FsClient) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

func (c *FsClient) CopyToRemote(src string, dst string) error {
	return nil
}

func (c *FsClient) Open(name string) (common.FsReader, error) {
	return nil, nil
}

func (c *FsClient) Close() error {
	return nil
}
