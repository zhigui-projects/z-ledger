package hdfs

import (
	"errors"
	"os"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/pengisgood/hdfs"
)

var logger = flogging.MustGetLogger("dfs.hdfs")

type FsClient struct {
	client *hdfs.Client
}

func NewFsClient(conf *common.HdfsConfig) (*FsClient, error) {
	if len(conf.NameNodes) == 0 {
		errMsg := "archive service can't be initialized, due to no namenode address in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if len(conf.User) == 0 {
		errMsg := "archive service can't be initialized, due to no HDFS user in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	options := hdfs.ClientOptions{
		Addresses:           conf.NameNodes,
		User:                conf.User,
		UseDatanodeHostname: conf.UseDatanodeHostname,
	}
	client, err := hdfs.NewClient(options)

	logger.Infof("Created a dfs client with options: %+v, error: %+v", options, err)
	return &FsClient{client: client}, err
}

func (c *FsClient) ReadDir(dirname string) ([]os.FileInfo, error) {
	return c.client.ReadDir(dirname)
}

func (c *FsClient) Stat(name string) (os.FileInfo, error) {
	return c.client.Stat(name)
}

func (c *FsClient) CopyToRemote(src string, dst string) error {
	return c.client.CopyToRemote(src, dst)
}

func (c *FsClient) Open(name string) (common.FsReader, error) {
	fileReader, err := c.client.Open(name)
	return &FsReader{fileReader: fileReader}, err
}

func (c *FsClient) Close() error {
	return c.client.Close()
}

type FsReader struct {
	fileReader *hdfs.FileReader
}

func (r *FsReader) Stat() os.FileInfo {
	return r.fileReader.Stat()
}

func (r *FsReader) Seek(offset int64, whence int) (int64, error) {
	return r.fileReader.Seek(offset, whence)
}

func (r *FsReader) Read(b []byte) (n int, err error) {
	return r.fileReader.Read(b)
}

func (r *FsReader) ReadAt(b []byte, offset int64) (int, error) {
	return r.fileReader.ReadAt(b, offset)
}

func (r *FsReader) Close() error {
	return r.fileReader.Close()
}
