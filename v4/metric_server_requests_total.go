package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelMethod,
		labelService,
	}
	return newRequestsTotalCounterVec("server", "requests_received_total", "TODO", labels, opts...)
}

type ServerRequestsTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewServerRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerRequestsTotalStatsHandler {
	h := &ServerRequestsTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options:   statsHandlerOptions{},
		},
		vec: vec,
	}
	h.baseStatsHandler.options.handleRPCLabelFn = h.serverRequestsTotalLabels
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

func (h *ServerRequestsTotalStatsHandler) serverRequestsTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	specialLabelValues := []string{
		tag.method,
		tag.service,
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
