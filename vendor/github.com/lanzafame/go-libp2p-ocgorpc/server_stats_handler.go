package ocgorpc

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-gorpc/stats"
	"go.opencensus.io/tag"
)

// statsTagRPC gets the metadata from gorpc context, extracts the encoded tags from
// it and creates a new tag.Map and puts them into the returned context.
func (h *ServerHandler) statsTagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	startTime := time.Now()
	if info == nil {
		logger.Infof("opencensus: TagRPC called with nil info.")
		return ctx
	}
	d := &rpcData{
		startTime: startTime,
		method:    info.FullMethodName,
	}
	propagated := h.extractPropagatedTags(ctx)
	ctx = tag.NewContext(ctx, propagated)
	ctx, _ = tag.New(ctx, tag.Upsert(KeyServerMethod, methodName(info.FullMethodName)))
	return context.WithValue(ctx, rpcDataKey, d)
}

// extractPropagatedTags creates a new tag map containing the tags extracted from the
// gorpc metadata.
func (h *ServerHandler) extractPropagatedTags(ctx context.Context) *tag.Map {
	buf, _ := ctx.Value(rpcDataKey).([]byte)
	if buf == nil {
		return nil
	}
	propagated, err := tag.Decode(buf)
	if err != nil {
		logger.Errorf("opencensus: Failed to decode tags from gorpc metadata failed to decode: %v", err)
		return nil
	}
	return propagated
}
