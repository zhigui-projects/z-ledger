package ocgorpc

import (
	"context"

	"github.com/libp2p/go-libp2p-gorpc/stats"

	"go.opencensus.io/trace"
)

// ServerHandler implements gorpc stats.Handler recording OpenCensus stats and
// traces. Use with gorpc servers.
//
// When installed (see Example), tracing metadata is read from inbound RPCs
// by default. If no tracing metadata is present, or if the tracing metadata is
// present but the SpanContext isn't sampled, then a new trace may be started
// (as determined by Sampler).
type ServerHandler struct {
	// StartOptions to use for to spans started around RPCs handled by this server.
	//
	// These will apply even if there is tracing metadata already
	// present on the inbound RPC but the SpanContext is not sampled. This
	// ensures that each service has some opportunity to be traced. If you would
	// like to not add any additional traces for this gorpc service, set:
	//
	//   StartOptions.Sampler = trace.ProbabilitySampler(0.0)
	//
	// StartOptions.SpanKind will always be set to trace.SpanKindServer
	// for spans started by this handler.
	StartOptions trace.StartOptions
}

var _ stats.Handler = (*ServerHandler)(nil)

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s *ServerHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	traceHandleRPC(ctx, rs)
	statsHandleRPC(ctx, rs)
}

// TagRPC implements per-RPC context management.
func (s *ServerHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	ctx = s.traceTagRPC(ctx, rti)
	ctx = s.statsTagRPC(ctx, rti)
	return ctx
}
