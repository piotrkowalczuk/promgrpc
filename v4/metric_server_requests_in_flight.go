package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	labels := []string{
		// keep alphabetical order
		labelMethod,
		labelService,
	}
	return newRequestsInFlightGaugeVec("server", labels, opts...)
}

type ServerRequestsInFlightStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewServerRequestsInFlightStatsHandler ...
func NewServerRequestsInFlightStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ServerRequestsInFlightStatsHandler {
	h := &ServerRequestsInFlightStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options:   statsHandlerOptions{},
		},
		vec: vec,
	}
	h.baseStatsHandler.options.handleRPCLabelFn = h.serverRequestsInFlightLabels
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

func (h *ServerRequestsInFlightStatsHandler) serverRequestsInFlightLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	// keep alphabetical order
	specialLabelValues := []string{
		tag.method,
		tag.service,
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
