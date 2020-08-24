package ipfs

import (
	"crypto/sha256"
	"io"
	"testing"

	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/stretchr/testify/assert"
)

var conf = &common.IpfsConfig{
	Url:         "139.224.117.14:5001",
	ClusterUrl:  "139.224.117.14:9094",
	CidIndexDir: "/Users/maxpeng/Projects/ipfs/cidIndex",
}

const filePath = "/Users/maxpeng/Projects/ipfs/hello.txt"

func TestSeek(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	reader, err := fsClient.Open(filePath)
	assert.Nil(t, err)

	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		logger.Errorf("io.Copy file[%s] got error: %s", filePath, err)
	}

	newPosition, err := reader.Seek(0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), newPosition)
}

func TestOpen(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	reader, err := fsClient.Open(filePath)
	assert.Nil(t, err)

	out := make([]byte, reader.Stat().Size())
	n, err := reader.Read(out)
	assert.Nil(t, err)
	assert.Equal(t, int(reader.Stat().Size()), n)
	assert.Equal(t, "hello from postman!\n", string(out))
}

func TestStat(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	r, err := fsClient.Open(filePath)
	assert.Nil(t, err)

	stat := r.Stat()
	assert.NotNil(t, stat)
	assert.Equal(t, filePath, stat.Name())
}

func TestCopyToRemote(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	err = fsClient.CopyToRemote(filePath, "")
	assert.Nil(t, err)
}

func TestReadDir(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	fileInfos, err := fsClient.ReadDir("/Users/maxpeng/Projects/ipfs")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(fileInfos))
}
