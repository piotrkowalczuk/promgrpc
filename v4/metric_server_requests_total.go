package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverRequestsTotalCounterVecSupportedLabels = supportedLabels{
	Method:  true,
	Service: true,
}

// NewServerRequestsTotalCounterVec instantiates default server-side CounterVec suitable for use with NewServerRequestsTotalStatsHandler.
func NewServerRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newRequestsTotalCounterVec("server", "requests_received_total", "TODO", serverRequestsTotalCounterVecSupportedLabels.labels(), opts...)
}

// ServerRequestsTotalStatsHandler dedicated server-side StatsHandlerCollector that counts number of requests sent.
type ServerRequestsTotalStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.CounterVec
}

// NewServerRequestsTotalStatsHandler instantiates ServerRequestsTotalStatsHandler based on given CounterVec and options.
// The CounterVec must have zero, one or two non-const non-curried labels.
// For those, the only allowed names are "grpc_method" and "grpc_service".
func NewServerRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerRequestsTotalStatsHandler {
	h := &ServerRequestsTotalStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverRequestsTotalCounterVecSupportedLabels),
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
func (h *ServerRequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if beg, ok := stat.(*stats.Begin); ok {
		switch {
		case !beg.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
