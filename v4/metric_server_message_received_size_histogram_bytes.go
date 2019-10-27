package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerMessageReceivedSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageReceivedSizeHistogramVec("server", opts...)
}

type ServerMessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewServerMessageReceivedSizeStatsHandler ...
func NewServerMessageReceivedSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerMessageReceivedSizeStatsHandler {
	h := &ServerMessageReceivedSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: messageReceivedSizeLabels,
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
func (h *ServerMessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.InPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
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
