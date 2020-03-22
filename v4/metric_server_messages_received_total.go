package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerMessagesReceivedTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelClientUserAgent,
		labelMethod,
		labelService,
	}
	return newMessagesReceivedTotalCounterVec("server", labels, opts...)
}

type ServerMessagesReceivedTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerMessagesReceivedTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service" and "user_agent".
func NewServerMessagesReceivedTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerMessagesReceivedTotalStatsHandler {
	h := &ServerMessagesReceivedTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: serverMessagesReceivedTotalLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ServerMessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.InPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func serverMessagesReceivedTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	return []string{
		tag.clientUserAgent,
		tag.method,
		tag.service,
	}
}
