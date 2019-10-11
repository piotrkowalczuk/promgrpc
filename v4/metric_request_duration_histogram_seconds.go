package promgrpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewRequestDurationHistogramVec(sub Subsystem, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "request_duration_histogram_seconds",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		[]string{
			labelClientUserAgent,
			labelCode,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

type RequestDurationStatsHandler struct {
	baseStatsHandler
	vec prometheus.ObserverVec
}

// NewRequestDurationStatsHandler ...
func NewRequestDurationStatsHandler(sub Subsystem, vec prometheus.ObserverVec, opts ...StatsHandlerOption) *RequestDurationStatsHandler {
	h := &RequestDurationStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
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
func (h *RequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}

func requestDurationLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
