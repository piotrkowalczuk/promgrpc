package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverMessagesSentTotalCounterVecSupportedLabels = supportedLabels{
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewServerMessagesSentTotalCounterVec instantiates default server-side CounterVec suitable for use with NewServerMessagesSentTotalStatsHandler.
func NewServerMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newMessagesSentTotalCounterVec("server", serverMessagesSentTotalCounterVecSupportedLabels.labels(), opts...)
}

// ServerMessagesSentTotalStatsHandler dedicated server-side StatsHandlerCollector that counts number of messages sent.
type ServerMessagesSentTotalStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewServerMessagesSentTotalStatsHandler instantiates ServerMessagesSentTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two or three non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewServerMessagesSentTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerMessagesSentTotalStatsHandler {
	h := &ServerMessagesSentTotalStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverMessagesSentTotalCounterVecSupportedLabels),
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
func (h *ServerMessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.OutPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
