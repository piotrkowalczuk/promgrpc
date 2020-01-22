package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	labels := []string{
		// keep alphabetical order
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newRequestsInFlightGaugeVec("client", labels, opts...)
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
				handleRPCLabelFn: clientRequestsInFlightLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

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

func clientRequestsInFlightLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	// keep alphabetical order
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
