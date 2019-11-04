package promgrpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	labels := []string{
		labelCode,
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newRequestDurationHistogramVec("client", labels, opts...)
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
				handleRPCLabelFn: clientRequestDurationLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

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

func clientRequestDurationLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
