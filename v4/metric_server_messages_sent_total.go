package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerMessagesSentTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelClientUserAgent,
		labelMethod,
		labelService,
	}
	return newMessagesSentTotalCounterVec("server", labels, opts...)
}

type ServerMessagesSentTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerMessagesSentTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewServerMessagesSentTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerMessagesSentTotalStatsHandler {
	h := &ServerMessagesSentTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: serverMessagesSentTotalLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ServerMessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.OutPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func serverMessagesSentTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	return []string{
		tag.clientUserAgent,
		tag.method,
		tag.service,
	}
}
