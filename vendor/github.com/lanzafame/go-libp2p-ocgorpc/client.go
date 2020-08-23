package ocgorpc

import (
	"context"

	"github.com/libp2p/go-libp2p-gorpc/stats"
	"go.opencensus.io/trace"
)

// ClientHandler implements a gorpc stats.Handler for recording OpenCensus stats and
// traces. Use with gorpc clients only.
type ClientHandler struct {
	// StartOptions allows configuring the StartOptions used to create new spans.
	//
	// StartOptions.SpanKind will always be set to trace.SpanKindClient
	// for spans started by this handler.
	StartOptions trace.StartOptions
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *ClientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	traceHandleRPC(ctx, rs)
	statsHandleRPC(ctx, rs)
}

// TagRPC implements per-RPC context management.
func (c *ClientHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	ctx = c.traceTagRPC(ctx, rti)
	ctx = c.statsTagRPC(ctx, rti)
	return ctx
}
