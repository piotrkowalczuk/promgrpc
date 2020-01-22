package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newMessagesSentTotalCounterVec("client", labels, opts...)
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
				handleRPCLabelFn: clientMessagesSentTotalLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

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

func clientMessagesSentTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
