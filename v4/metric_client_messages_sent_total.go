package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newMessagesSentTotalCounterVec("client", opts...)
}

type ClientMessagesSentTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewClientMessagesSentTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewClientMessagesSentTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientMessagesSentTotalStatsHandler {
	h := &ClientMessagesSentTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: messagesSentTotalLabels,
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
func (h *ClientMessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
