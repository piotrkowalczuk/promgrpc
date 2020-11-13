package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverMessagesReceivedTotalCounterVecSupportedLabels = supportedLabels{
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewServerMessagesReceivedTotalCounterVec instantiates default server-side CounterVec suitable for use with NewServerMessagesReceivedTotalStatsHandler.
func NewServerMessagesReceivedTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newMessagesReceivedTotalCounterVec("server", serverMessagesReceivedTotalCounterVecSupportedLabels.labels(), opts...)
}

// ServerMessagesReceivedTotalStatsHandler dedicated server-side StatsHandlerCollector that counts number of messages received.
type ServerMessagesReceivedTotalStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewServerMessagesReceivedTotalStatsHandler instantiates ServerMessagesReceivedTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two or three non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewServerMessagesReceivedTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerMessagesReceivedTotalStatsHandler {
	h := &ServerMessagesReceivedTotalStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverMessageSentSizeHistogramVecSupportedLabels),
	}
	h.baseStatsHandler = baseStatsHandler{
		collector: vec,
		options: statsHandlerOptions{
			handleRPCLabelFn: h.labelsTagRPC,
		},
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ServerMessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.InPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
