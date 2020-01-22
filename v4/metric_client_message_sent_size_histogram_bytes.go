package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientMessageSentSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	labels := []string{
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newMessageSentSizeHistogramVec("client", labels, opts...)
}

type ClientMessageSentSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewMessageSentSizeStatsHandler ...
func NewClientMessageSentSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientMessageSentSizeStatsHandler {
	h := &ClientMessageSentSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: clientMessageSentSizeLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ClientMessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}

func clientMessageSentSizeLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
