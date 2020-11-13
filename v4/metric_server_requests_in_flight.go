package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverRequestsInFlightGaugeVecSupportedLabels = supportedLabels{
	Method:  true,
	Service: true,
}

// NewServerRequestsInFlightGaugeVec instantiates default server-side GaugeVec suitable for use with NewServerRequestsInFlightStatsHandler.
func NewServerRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newRequestsInFlightGaugeVec("server", serverRequestsInFlightGaugeVecSupportedLabels.labels(), opts...)
}

// ServerRequestsInFlightStatsHandler dedicated server-side StatsHandlerCollector that counts the number of requests currently in flight.
type ServerRequestsInFlightStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.GaugeVec
}

// NewServerRequestsInFlightStatsHandler instantiates ServerRequestsInFlightStatsHandler based on given GaugeVec and options.
// The GaugeVec must have zero, one or two non-const non-curried labels.
// For those, the only allowed names are "grpc_method" and "grpc_service".
func NewServerRequestsInFlightStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ServerRequestsInFlightStatsHandler {
	h := &ServerRequestsInFlightStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverRequestsInFlightGaugeVecSupportedLabels),
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

// HandleRPC processes the RPC stats.
func (h *ServerRequestsInFlightStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.Begin:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.End:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Dec()
		}
	}
}
