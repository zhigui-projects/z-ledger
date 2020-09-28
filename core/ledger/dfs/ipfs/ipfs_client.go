package ipfs

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	sh "github.com/pengisgood/go-ipfs-api"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var logger = flogging.MustGetLogger("dfs.ipfs")

type FsClient struct {
	shell         *sh.Shell
	clusterClient *ClusterClient
	cidStore      *leveldb.DB
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

// Sys returns the cid from IPFS
func (fi *fileInfo) Cid() string {
	return fi.hash
}

func NewFsClient(conf *common.IpfsConfig) (*FsClient, error) {
	if len(conf.Url) == 0 {
		errMsg := "archive service can't be initialized, due to no IPFS url in configuration"
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	shell := sh.NewShell(conf.Url)
	clusterClient := NewClusterClient(conf.ClusterUrl)

	err := createDirIfMissing(conf.CidIndexDir)
	if err != nil {
		return nil, err
	}
	db, err := leveldb.OpenFile(conf.CidIndexDir, nil)
	if err != nil {
		logger.Debugf("Error creating level db [%s]", conf.CidIndexDir)
		return nil, errors.New(fmt.Sprintf("creating dir [%s] with err: %s", conf.CidIndexDir, err))
	}

	logger.Infof("Created a dfs client with options: %+v", conf)
	return &FsClient{shell: shell, clusterClient: clusterClient, cidStore: db}, nil
}

func createDirIfMissing(dirPath string) error {
	// if dirPath does not end with a path separator, it leaves out the last segment while creating directories
	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}
	logger.Debugf("CreateDirIfMissing [%s]", dirPath)
	err := os.MkdirAll(path.Dir(dirPath), 0755)
	if err != nil {
		logger.Debugf("Error creating dir [%s]", dirPath)
		return errors.New(fmt.Sprintf("error creating dir [%s] with err: %s", dirPath, err))
	}
	return nil
}

type cidInfo struct {
	Cid  string
	Name string
	Size int64
}

func encode(value *cidInfo) []byte {
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(value); err != nil {
		return nil
	}
	return buf.Bytes()
}

func decode(value []byte) cidInfo {
	var out cidInfo
	decoder := gob.NewDecoder(bytes.NewBuffer(value))
	if err := decoder.Decode(&out); err != nil {
		return cidInfo{}
	}
	return out
}

func (c *FsClient) ReadDir(dirname string) ([]os.FileInfo, error) {
	iter := c.cidStore.NewIterator(util.BytesPrefix([]byte(dirname)), nil)
	var infos []os.FileInfo
	for iter.Next() {
		value := decode(iter.Value())
		infos = append(infos, &fileInfo{value.Name, 0, value.Size, value.Cid})
	}
	iter.Release()
	return infos, nil
}

func (c *FsClient) Stat(name string) (os.FileInfo, error) {
	value, err := c.cidStore.Get([]byte(name), nil)
	if err != nil {
		logger.Errorf("stat file[%s], got error[%s]", name, err)
		return nil, err
	}
	cidInfo := decode(value)
	logger.Infof("stat file[%s], got value[%s]", name, cidInfo)
	return &fileInfo{cidInfo.Name, 0, cidInfo.Size, cidInfo.Cid}, nil
}

func (c *FsClient) CopyToRemote(src string, dst string) error {
	var file *os.File
	var err error
	if file, err = os.Open(src); err != nil {
		logger.Errorf("open file[%s], got error[%s]", src, err)
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		logger.Errorf("stat file[%s], got error[%s]", src, err)
		return err
	}

	cid, err := c.clusterClient.Add(file)
	if err != nil {
		logger.Errorf("add file[%s] to ipfs cluster, got error[%s]", src, err)
		return err
	}

	value := &cidInfo{Cid: cid, Name: dst, Size: stat.Size()}
	logger.Infof("cidStore PUT[key: %s, value: %s]", dst, value)
	if err := c.cidStore.Put([]byte(dst), encode(value), nil); err != nil {
		logger.Errorf("add file[%s], cid[%s] to cid index store, got error[%s]", src, cid, err)
		return err
	}

	return nil
}

func (c *FsClient) Open(name string) (common.FsReader, error) {
	stat, err := c.Stat(name)
	if err != nil {
		logger.Errorf("stat file[%s], got error: %s", name, err)
		return nil, err
	}

	cid := stat.(*fileInfo).Cid()

	resp, err := c.shell.Request("cat", cid).
		Option("length", stat.Size()).
		Send(context.Background())
	if err != nil {
		logger.Errorf("cat file[%s] from ipfs, got error: %s", cid, err)
		return nil, err
	}
	if resp.Error != nil {
		logger.Errorf("cat file[%s] from ipfs, got resp.error: %s", cid, resp.Error)
		return nil, resp.Error
	}

	output := make([]byte, stat.Size())
	var actualLength int
	if actualLength, err = resp.Output.Read(output); err != nil && err.Error() != "EOF" {
		logger.Errorf("reader.Read with stat[%s], got error: %s", stat, err)
		return nil, err
	}

	logger.Infof("read from ipfs, expected length: %d, actual length: %d", stat.Size(), actualLength)

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
