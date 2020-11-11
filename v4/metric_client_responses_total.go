package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientResponsesTotalCounterVecSupportedLabels = supportedLabels{
	Code:            true,
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientResponsesTotalCounterVec instantiates client-side CounterVec suitable for use with NewClientResponsesTotalStatsHandler.
func NewClientResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newResponsesTotalCounterVec("client", "responses_received_total", "TODO", clientResponsesTotalCounterVecSupportedLabels.labels(), opts...)
}

// ClientResponsesTotalStatsHandler dedicated client-side StatsHandlerCollector that counts number of responses received.
type ClientResponsesTotalStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewClientResponsesTotalStatsHandler instantiates ClientResponsesTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two, three, four or five non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service", "grpc_client_user_agent" and "grpc_code".
func NewClientResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientResponsesTotalStatsHandler {
	h := &ClientResponsesTotalStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientResponsesTotalCounterVecSupportedLabels),
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
func (h *ClientResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.End:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
