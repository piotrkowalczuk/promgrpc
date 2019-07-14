package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewMessagesReceivedTotalCounterVec(sub Subsystem, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "messages_received_total",
		Help:      "TODO",
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

// MessagesReceivedTotalLabels ...
func MessagesReceivedTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}

type MessagesReceivedTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewMessagesReceivedTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service" and "user_agent".
func NewMessagesReceivedTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec, opts ...StatsHandlerOption) *MessagesReceivedTotalStatsHandler {
	h := &MessagesReceivedTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				rpcLabelFn: MessagesReceivedTotalLabels,
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
func (h *MessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.InPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		}
	}
}
