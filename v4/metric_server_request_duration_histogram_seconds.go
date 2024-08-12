package promgrpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	labels := []string{
		labelClientUserAgent,
		labelCode,
		labelMethod,
		labelService,
	}

	return newRequestDurationHistogramVec("server", labels, opts...)
}

type ServerRequestDurationStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewServerRequestDurationStatsHandler ...
func NewServerRequestDurationStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerRequestDurationStatsHandler {
	h := &ServerRequestDurationStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options:   statsHandlerOptions{},
		},
		vec: vec,
	}
	h.baseStatsHandler.options.handleRPCLabelFn = h.serverRequestDurationLabels
	h.applyOpts(opts...)

	return h
}

// HandleRPC processes the RPC stats.
func (h *ServerRequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		switch {
		case !stat.IsClient():
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}

func (h *ServerRequestDurationStatsHandler) serverRequestDurationLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	specialLabelValues := []string{
		tag.clientUserAgent,
		status.Code(stat.(*stats.End).Error).String(),
		tag.method,
		tag.service,
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}

	return specialLabelValues
}
