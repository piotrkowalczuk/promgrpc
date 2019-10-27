package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newRequestDurationHistogramVec("client", opts...)
}

type ClientRequestDurationStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewClientRequestDurationStatsHandler ...
func NewClientRequestDurationStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientRequestDurationStatsHandler {
	h := &ClientRequestDurationStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: requestDurationLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC processes the RPC stats.
func (h *ClientRequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		switch {
		case stat.IsClient():
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}
