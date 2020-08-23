package ocgorpc

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-gorpc/stats"
	ocstats "go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var logger = logging.Logger("p2p-gorpc")

type gorpcInstrumentationKey string

// rpcData holds the instrumentation RPC data that is needed between the start
// and end of an call. It holds the info that this package needs to keep track
// of between the various gorpc events.
type rpcData struct {
	// reqCount and respCount has to be the first words
	// in order to be 64-aligned on 32-bit architectures.
	sentCount, sentBytes, recvCount, recvBytes int64 // access atomically

	// startTime represents the time at which TagRPC was invoked at the
	// beginning of an RPC. It is an appoximation of the time when the
	// application code invoked gorpc code.
	startTime time.Time
	method    string
}

// The following variables define the default hard-coded auxiliary data used by
// both the default gorpc client and gorpc server metrics.
var (
	DefaultBytesDistribution        = view.Distribution(0, 1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296)
	DefaultMillisecondsDistribution = view.Distribution(0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
	DefaultMessageCountDistribution = view.Distribution(0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536)
)

// Server tags are applied to the context used to process each RPC, as well as
// the measures at the end of each RPC.
var (
	KeyServerMethod, _ = tag.NewKey("gorpc_server_method")
	KeyServerStatus, _ = tag.NewKey("gorpc_server_status")
)

// Client tags are applied to measures at the end of each RPC.
var (
	KeyClientMethod, _ = tag.NewKey("gorpc_client_method")
	KeyClientStatus, _ = tag.NewKey("gorpc_client_status")
)

var (
	rpcDataKey = gorpcInstrumentationKey("opencensus-rpcData")
)

func methodName(fullname string) string {
	return strings.TrimLeft(fullname, "/")
}

// statsHandleRPC processes the RPC events.
func statsHandleRPC(ctx context.Context, s stats.RPCStats) {
	switch st := s.(type) {
	case *stats.Begin:
		// do nothing for client
	case *stats.OutPayload:
		handleRPCOutPayload(ctx, st)
	case *stats.InPayload:
		handleRPCInPayload(ctx, st)
	case *stats.End:
		handleRPCEnd(ctx, st)
	default:
		logger.Infof("unexpected stats: %T", st)
	}
}

func handleRPCOutPayload(ctx context.Context, s *stats.OutPayload) {
	d, ok := ctx.Value(rpcDataKey).(*rpcData)
	if !ok {
		logger.Info("Failed to retrieve *rpcData from context.")
		return
	}

	atomic.AddInt64(&d.sentBytes, int64(s.Length))
	atomic.AddInt64(&d.sentCount, 1)
}

func handleRPCInPayload(ctx context.Context, s *stats.InPayload) {
	d, ok := ctx.Value(rpcDataKey).(*rpcData)
	if !ok {
		logger.Info("Failed to retrieve *rpcData from context.")
		return
	}

	atomic.AddInt64(&d.recvBytes, int64(s.Length))
	atomic.AddInt64(&d.recvCount, 1)
}

func handleRPCEnd(ctx context.Context, s *stats.End) {
	d, ok := ctx.Value(rpcDataKey).(*rpcData)
	if !ok {
		logger.Info("Failed to retrieve *rpcData from context.")
		return
	}

	elapsedTime := time.Since(d.startTime)

	// TODO(lanzafame): add status code-esque concept to gorpc
	// var st string
	// if s.Error != nil {
	// } else {
	// 	st = "OK"
	// }

	latencyMillis := float64(elapsedTime) / float64(time.Millisecond)
	if s.Client {
		ocstats.RecordWithTags(
			ctx,
			[]tag.Mutator{
				tag.Upsert(KeyClientMethod, methodName(d.method)),
				// tag.Upsert(KeyClientStatus, st),
			},
			ClientSentBytesPerRPC.M(atomic.LoadInt64(&d.sentBytes)),
			ClientSentMessagesPerRPC.M(atomic.LoadInt64(&d.sentCount)),
			ClientReceivedMessagesPerRPC.M(atomic.LoadInt64(&d.recvCount)),
			ClientReceivedBytesPerRPC.M(atomic.LoadInt64(&d.recvBytes)),
			ClientRoundtripLatency.M(latencyMillis),
		)
	} else {
		ocstats.RecordWithTags(
			ctx,
			[]tag.Mutator{
				// tag.Upsert(KeyServerStatus, st),
			},
			ServerSentBytesPerRPC.M(atomic.LoadInt64(&d.sentBytes)),
			ServerSentMessagesPerRPC.M(atomic.LoadInt64(&d.sentCount)),
			ServerReceivedMessagesPerRPC.M(atomic.LoadInt64(&d.recvCount)),
			ServerReceivedBytesPerRPC.M(atomic.LoadInt64(&d.recvBytes)),
			ServerLatency.M(latencyMillis),
		)
	}
}
