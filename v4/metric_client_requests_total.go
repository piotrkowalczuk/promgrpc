package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newRequestsTotalCounterVec("client", "requests_sent_total", "TODO", opts...)
}

type ClientRequestsTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewClientRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewClientRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientRequestsTotalStatsHandler {
	h := &ClientRequestsTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: requestsTotalLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC implements stats Handler interface.
func (h *ClientRequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if !stat.IsClient() {
		return
	}
	if _, ok := stat.(*stats.Begin); ok {
		h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
	}
}
