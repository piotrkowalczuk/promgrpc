package promgrpc

import (
	"context"

	"github.com/piotrkowalczuk/promgrpc/v4/internal/useragent"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelIsFailFast,
		labelMethod,
		labelService,
		labelClientUserAgent,
	}
	return newMessagesSentTotalCounterVec("client", labels, opts...)
}

type ClientMessagesSentTotalStatsHandler struct {
	baseStatsHandler
	uas useragent.Store
	vec *prometheus.CounterVec
}

// NewClientMessagesSentTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewClientMessagesSentTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientMessagesSentTotalStatsHandler {
	h := &ClientMessagesSentTotalStatsHandler{
		vec: vec,
	}
	h.baseStatsHandler = baseStatsHandler{
		collector: vec,
		options: statsHandlerOptions{
			handleRPCLabelFn: h.labels,
		},
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

func (h *ClientMessagesSentTotalStatsHandler) labels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
		h.uas.ClientSide(ctx, stat),
	}
}
