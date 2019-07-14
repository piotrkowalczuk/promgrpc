package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewMessageSentSizeHistogramVec(sub Subsystem, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "message_sent_size_histogram_bytes",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

type MessageSentSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewMessageSentSizeStatsHandler ...
func NewMessageSentSizeStatsHandler(sub Subsystem, vec prometheus.ObserverVec, opts ...StatsHandlerOption) *MessageSentSizeStatsHandler {
	h := &MessageSentSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				rpcLabelFn: messageSentSizeLabels,
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
func (h *MessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}

func messageSentSizeLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
