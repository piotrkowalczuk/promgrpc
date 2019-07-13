package promgrpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewRequestDurationHistogramVec(sub Subsystem) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "request_duration_histogram_seconds",
			Help:      "TODO",
		},
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
	vec *prometheus.HistogramVec
}

// NewRequestDurationStatsHandler ...
func NewRequestDurationStatsHandler(sub Subsystem, vec *prometheus.HistogramVec) *RequestDurationStatsHandler {
	return &RequestDurationStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// HandleRPC processes the RPC stats.
func (h *RequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		tag := ctx.Value(tagRPCKey).(rpcTag)
		lab := []string{
			tag.clientUserAgent,
			status.Code(end.Error).String(),
			tag.isFailFast,
			tag.method,
			tag.service,
		}
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(lab...).Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(lab...).Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}
