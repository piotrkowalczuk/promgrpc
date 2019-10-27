package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerMessageSentSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageSentSizeHistogramVec("server", opts...)
}

type ServerMessageSentSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewServerMessageSentSizeStatsHandler ...
func NewServerMessageSentSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerMessageSentSizeStatsHandler {
	h := &ServerMessageSentSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: messageSentSizeLabels,
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
func (h *ServerMessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.OutPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
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
