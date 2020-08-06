package ipfs

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	sh "github.com/pengisgood/go-ipfs-api"
)

var logger = flogging.MustGetLogger("dfs.ipfs")

type FsClient struct {
	shell *sh.Shell
}

type fileInfo struct {
	name  string
	_type uint8
	size  int64
	hash  string
}

func (fi *fileInfo) Name() string {
	return fi.name
}
func (fi *fileInfo) Size() int64 {
	return fi.size
}
func (fi *fileInfo) Mode() os.FileMode {
	return 0444
}
func (fi *fileInfo) ModTime() time.Time {
	return time.Unix(0, 0)
}
func (fi *fileInfo) IsDir() bool {
	switch fi._type {
	case 0:
		return false
	case 1:
		return true
	default:
		logger.Warnf("Unknown IPFS file type[%s]", fi._type)
		return false
	}
}

// Sys returns the hash from IPFS
func (fi *fileInfo) Sys() interface{} {
	return fi.hash
}

func NewFsClient(conf *common.IpfsConfig) (*FsClient, error) {
	if len(conf.Url) == 0 {
		errMsg := "archive service can't be initialized, due to no IPFS url in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	shell := sh.NewShell(conf.Url)
	logger.Infof("Created a dfs client with options: %+v", conf)
	return &FsClient{shell: shell}, nil
}

func (c *FsClient) ReadDir(dirname string) ([]os.FileInfo, error) {
	entries, err := c.shell.FilesLs(dirname)
	if err != nil {
		logger.Errorf("read dir[%s] from MFS, got error[%s]", dirname, err)
		return nil, err
	}

	var infos []os.FileInfo
	for _, entry := range entries {
		infos = append(infos, &fileInfo{entry.Name, entry.Type, int64(entry.Size), entry.Hash})
	}
	return infos, nil
}

func (c *FsClient) Stat(name string) (os.FileInfo, error) {
	stat, err := c.shell.FilesStat(name)
	if err != nil {
		logger.Errorf("stat file[%s] from MFS, got error[%s]", name, err)
		return nil, err
	}

	return &fileInfo{name, 0, int64(stat.Size), stat.Hash}, nil
}

func (c *FsClient) CopyToRemote(src string, dst string) error {
	var file *os.File
	var err error
	if file, err = os.Open(src); err != nil {
		logger.Errorf("copy file[%s] to ipfs, got error[%s]", src, err)
		return err
	}
	if err = c.shell.FilesWrite(dst, file); err != nil {
		logger.Errorf("copy file[%s] to ipfs, got error[%s]", src, err)
		return err
	}

	logger.Infof("copy file[%s] to MFS", src)
	return nil
}

func (c *FsClient) Open(name string) (common.FsReader, error) {
	stat, err := c.Stat(name)
	if err != nil {
		logger.Errorf("stat file[%s], got error: %s", name, err)
		return nil, err
	}

	r, err := c.shell.FilesRead(name)
	if err != nil {
		logger.Errorf("read file[%s] from ipfs, got error: %s", name, err)
		return nil, err
	}

	output := make([]byte, stat.Size())
	if _, err := r.Read(output); err != nil && err.Error() != "EOF" {
		logger.Errorf("reader.Read with stat[%s], got error: %s", stat, err)
		return nil, err
	}

	return &FsReader{reader: bytes.NewReader(output), stat: stat}, nil
}

type FsReader struct {
	reader *bytes.Reader
	stat   os.FileInfo
}

func (r *FsReader) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}

func (r *FsReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *FsReader) Close() error {
	// do nothing
	return nil
}

func (r *FsReader) Stat() os.FileInfo {
	return r.stat
}

func (r *FsReader) ReadAt(b []byte, offset int64) (int, error) {
	return r.reader.ReadAt(b, offset)
}

func (c *FsClient) Close() error {
	// do nothing
	return nil
}
