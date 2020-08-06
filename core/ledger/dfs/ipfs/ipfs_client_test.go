package ipfs

import (
	"testing"

	"github.com/hyperledger/fabric/core/ledger/dfs/common"
	"github.com/stretchr/testify/assert"
)

func TestSeek(t *testing.T) {
	conf := &common.IpfsConfig{Url: "47.52.139.52:5001"}
	fsClient, err := NewFsClient(conf)
	assert.Nil(t, err)

	r, err := fsClient.Open("/blk/var/hyperledger/production/ledgersData/chains/chains/mychannel/blockfile_000002")
	assert.Nil(t, err)

	newPosition, err := r.Seek(0, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), newPosition)
}
