package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newRequestsTotalCounterVec("client", "requests_sent_total", "TODO", labels, opts...)
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
				handleRPCLabelFn: clientRequestsTotalLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

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

func clientRequestsTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
