package eventbus

import (
	"github.com/asaskevich/EventBus"
	"github.com/hyperledger/fabric/common/flogging"
	"sync"
)

var (
	logger = flogging.MustGetLogger("archive.eventbus")
	bus    = make(map[string]EventBus.Bus)

	lock sync.RWMutex
)

const (
	ArchiveByTxDate = "ARCHIVE_BY_TX_DATE"
)

// Get returns the event bus singleton
func Get(channelId string) EventBus.Bus {
	lock.Lock()
	defer lock.Unlock()

	v, ok := bus[channelId]

	if !ok {
		logger.Infof("Initializing eventbus for channel[%s]", channelId)
		v = EventBus.New()
		bus[channelId] = v
	}
	return v
}
