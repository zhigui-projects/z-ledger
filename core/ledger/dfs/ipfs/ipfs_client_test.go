package ipfs

import (
	"crypto/sha256"
	"io"
	"testing"

	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/stretchr/testify/assert"
)

var conf = &common.IpfsConfig{Url: "47.52.139.52:5001"}

const filePath = "/blk/var/hyperledger/production/ledgersData/chains/chains/mychannel/blockfile_000002"

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

func TestStat(t *testing.T) {
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	r, err := fsClient.Open(filePath)
	assert.Nil(t, err)

	stat := r.Stat()
	assert.NotNil(t, stat)
	assert.Equal(t, filePath, stat.Name())
}
