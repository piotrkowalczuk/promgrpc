package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverResponsesTotalCounterVecSupportedLabels = supportedLabels{
	Method:          true,
	Service:         true,
	Code:            true,
	ClientUserAgent: true,
}

// NewServerResponsesTotalCounterVec instantiates default server-side CounterVec suitable for use with NewServerResponsesTotalStatsHandler.
func NewServerResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newResponsesTotalCounterVec("server", "responses_sent_total", "TODO", serverResponsesTotalCounterVecSupportedLabels.labels(), opts...)
}

// ServerResponsesTotalStatsHandler dedicated server-side StatsHandlerCollector that counts number of responses received.
type ServerResponsesTotalStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewServerResponsesTotalStatsHandler instantiates ServerResponsesTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service", "grpc_client_user_agent" and "grpc_code".
func NewServerResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerResponsesTotalStatsHandler {
	h := &ServerResponsesTotalStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverResponsesTotalCounterVecSupportedLabels),
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
func (h *ServerResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.End); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
