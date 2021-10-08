package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientMessagesReceivedTotalCounterVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientMessagesReceivedTotalCounterVec instantiates client-side CounterVec suitable for use with NewClientMessagesReceivedTotalStatsHandler.
func NewClientMessagesReceivedTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newMessagesReceivedTotalCounterVec("client", clientMessagesReceivedTotalCounterVecSupportedLabels.labels(), opts...)
}

// ClientMessagesReceivedTotalStatsHandler dedicated client-side StatsHandlerCollector that counts number of messages received.
type ClientMessagesReceivedTotalStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewClientMessagesReceivedTotalStatsHandler instantiates ClientMessagesReceivedTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientMessagesReceivedTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientMessagesReceivedTotalStatsHandler {
	h := &ClientMessagesReceivedTotalStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientMessagesReceivedTotalCounterVecSupportedLabels),
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
func (h *ClientMessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.InPayload:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
