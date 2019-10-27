package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientMessageReceivedSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageReceivedSizeHistogramVec("client", opts...)
}

type ClientMessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewMessageReceivedSizeStatsHandler ...
func NewClientMessageReceivedSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientMessageReceivedSizeStatsHandler {
	h := &ClientMessageReceivedSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: messageReceivedSizeLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ClientMessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.InPayload); ok {
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}
