package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newRequestsInFlightGaugeVec("client", opts...)
}

type ClientRequestsInFlightStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewClientRequestsInFlightStatsHandler ...
func NewClientRequestsInFlightStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ClientRequestsInFlightStatsHandler {
	h := &ClientRequestsInFlightStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: requestsInFlightLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC processes the RPC stats.
func (h *ClientRequestsInFlightStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.Begin:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.End:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Dec()
		}
	}
}
