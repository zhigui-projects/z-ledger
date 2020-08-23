package ocgorpc

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following variables are measures are recorded by ServerHandler:
var (
	ServerReceivedMessagesPerRPC = stats.Int64("gorpc.libp2p.io/server/received_messages_per_rpc", "Number of messages received in each RPC. Has value 1 for non-streaming RPCs.", stats.UnitDimensionless)
	ServerReceivedBytesPerRPC    = stats.Int64("gorpc.libp2p.io/server/received_bytes_per_rpc", "Total bytes received across all messages per RPC.", stats.UnitBytes)
	ServerSentMessagesPerRPC     = stats.Int64("gorpc.libp2p.io/server/sent_messages_per_rpc", "Number of messages sent in each RPC. Has value 1 for non-streaming RPCs.", stats.UnitDimensionless)
	ServerSentBytesPerRPC        = stats.Int64("gorpc.libp2p.io/server/sent_bytes_per_rpc", "Total bytes sent in across all response messages per RPC.", stats.UnitBytes)
	ServerLatency                = stats.Float64("gorpc.libp2p.io/server/server_latency", "Time between first byte of request received to last byte of response sent, or terminal error.", stats.UnitMilliseconds)
)

// TODO(acetechnologist): This is temporary and will need to be replaced by a
// mechanism to load these defaults from a common repository/config shared by
// all supported languages. Likely a serialized protobuf of these defaults.

// Predefined views may be registered to collect data for the above measures.
// As always, you may also define your own custom views over measures collected by this
// package. These are declared as a convenience only; none are registered by
// default.
var (
	ServerReceivedBytesPerRPCView = &view.View{
		Name:        "gorpc.libp2p.io/server/received_bytes_per_rpc",
		Description: "Distribution of received bytes per RPC, by method.",
		Measure:     ServerReceivedBytesPerRPC,
		TagKeys:     []tag.Key{KeyServerMethod},
		Aggregation: DefaultBytesDistribution,
	}

	ServerSentBytesPerRPCView = &view.View{
		Name:        "gorpc.libp2p.io/server/sent_bytes_per_rpc",
		Description: "Distribution of total sent bytes per RPC, by method.",
		Measure:     ServerSentBytesPerRPC,
		TagKeys:     []tag.Key{KeyServerMethod},
		Aggregation: DefaultBytesDistribution,
	}

	ServerLatencyView = &view.View{
		Name:        "gorpc.libp2p.io/server/server_latency",
		Description: "Distribution of server latency in milliseconds, by method.",
		TagKeys:     []tag.Key{KeyServerMethod},
		Measure:     ServerLatency,
		Aggregation: DefaultMillisecondsDistribution,
	}

	ServerCompletedRPCsView = &view.View{
		Name:        "gorpc.libp2p.io/server/completed_rpcs",
		Description: "Count of RPCs by method and status.",
		// TagKeys:     []tag.Key{KeyServerMethod, KeyServerStatus},
		TagKeys:     []tag.Key{KeyServerMethod},
		Measure:     ServerLatency,
		Aggregation: view.Count(),
	}

	ServerReceivedMessagesPerRPCView = &view.View{
		Name:        "gorpc.libp2p.io/server/received_messages_per_rpc",
		Description: "Distribution of messages received count per RPC, by method.",
		TagKeys:     []tag.Key{KeyServerMethod},
		Measure:     ServerReceivedMessagesPerRPC,
		Aggregation: DefaultMessageCountDistribution,
	}

	ServerSentMessagesPerRPCView = &view.View{
		Name:        "gorpc.libp2p.io/server/sent_messages_per_rpc",
		Description: "Distribution of messages sent count per RPC, by method.",
		TagKeys:     []tag.Key{KeyServerMethod},
		Measure:     ServerSentMessagesPerRPC,
		Aggregation: DefaultMessageCountDistribution,
	}
)

// DefaultServerViews are the default server views provided by this package.
var DefaultServerViews = []*view.View{
	ServerReceivedBytesPerRPCView,
	ServerSentBytesPerRPCView,
	ServerLatencyView,
	ServerCompletedRPCsView,
}
