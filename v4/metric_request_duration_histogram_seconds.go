package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
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
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
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

// Init implements StatsHandlerCollector interface.
func (h *RequestDurationStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC processes the RPC stats.
func (h *RequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	switch s := stat.(type) {
	case *stats.End:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Observe(s.EndTime.Sub(s.BeginTime).Seconds())
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Observe(s.EndTime.Sub(s.BeginTime).Seconds())
		}
	}
}
