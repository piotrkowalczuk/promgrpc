package promgrpc

import (
	"context"
	"strconv"
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
		[]string{labelIsFailFast, labelService, labelMethod, labelCode, labelClientUserAgent},
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
		lab := prometheus.Labels{
			labelMethod:          tag.method,
			labelService:         tag.service,
			labelIsFailFast:      strconv.FormatBool(tag.isFailFast),
			labelCode:            status.Code(end.Error).String(),
			labelClientUserAgent: tag.clientUserAgent,
		}
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}
