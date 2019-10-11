package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewMessagesSentTotalCounterVec(sub Subsystem, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "messages_sent_total",
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

type MessagesSentTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewMessagesSentTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewMessagesSentTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec, opts ...StatsHandlerOption) *MessagesSentTotalStatsHandler {
	h := &MessagesSentTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
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
func (h *MessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func messagesSentTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
