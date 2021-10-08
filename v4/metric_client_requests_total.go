package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientRequestsTotalCounterVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientRequestsTotalCounterVec instantiates client-side CounterVec suitable for use with NewClientRequestsTotalStatsHandler.
func NewClientRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newRequestsTotalCounterVec("client", "requests_sent_total", "TODO", clientRequestsTotalCounterVecSupportedLabels.labels(), opts...)
}

// ClientRequestsTotalStatsHandler dedicated client-side StatsHandlerCollector that counts number of requests sent.
type ClientRequestsTotalStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewClientRequestsTotalStatsHandler instantiates ClientRequestsTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientRequestsTotalStatsHandler {
	h := &ClientRequestsTotalStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientRequestsTotalCounterVecSupportedLabels),
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
func (h *ClientRequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if !stat.IsClient() {
		return
	}
	if _, ok := stat.(*stats.OutHeader); ok {
		h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
	}
}
