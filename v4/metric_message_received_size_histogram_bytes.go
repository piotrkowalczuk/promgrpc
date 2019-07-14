package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewMessageReceivedSizeHistogramVec(sub Subsystem, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "message_received_size_histogram_bytes",
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

type MessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewMessageReceivedSizeStatsHandler ...
func NewMessageReceivedSizeStatsHandler(sub Subsystem, vec prometheus.ObserverVec, opts ...StatsHandlerOption) *MessageReceivedSizeStatsHandler {
	h := &MessageReceivedSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				rpcLabelFn: messageReceivedSizeLabels,
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
func (h *MessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.InPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}

func messageReceivedSizeLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
