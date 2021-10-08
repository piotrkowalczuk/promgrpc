package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientMessagesSentTotalCounterVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientMessagesSentTotalCounterVec instantiates client-side CounterVec suitable for use with NewClientMessagesSentTotalStatsHandler.
func NewClientMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newMessagesSentTotalCounterVec("client", clientMessagesSentTotalCounterVecSupportedLabels.labels(), opts...)
}

// ClientMessagesSentTotalStatsHandler dedicated client-side StatsHandlerCollector that counts number of messages sent.
type ClientMessagesSentTotalStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewClientMessagesSentTotalStatsHandler instantiates ClientMessagesSentTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientMessagesSentTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientMessagesSentTotalStatsHandler {
	h := &ClientMessagesSentTotalStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientMessagesSentTotalCounterVecSupportedLabels),
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
func (h *ClientMessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.OutPayload:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
