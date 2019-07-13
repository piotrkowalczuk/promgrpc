package promgrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewRequestsGaugeVec(sub Subsystem) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "requests_in_flight",
			Help:      "TODO",
		},
		[]string{labelIsFailFast, labelService, labelMethod},
	)
}

type RequestsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
	idx map[rpcTag]prometheus.Gauge
}

// NewRequestsStatsHandler ...
func NewRequestsStatsHandler(sub Subsystem, vec *prometheus.GaugeVec) *RequestsStatsHandler {
	return &RequestsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
		idx: make(map[rpcTag]prometheus.Gauge),
	}
}

// HandleRPC processes the RPC stats.
func (h *RequestsStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.Begin:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(h.labels(ctx)).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(h.labels(ctx)).Inc()
		}
	case *stats.End:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(h.labels(ctx)).Dec()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(h.labels(ctx)).Dec()
		}
	}
}

func (h *RequestsStatsHandler) labels(ctx context.Context) prometheus.Labels {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return  prometheus.Labels{
		labelMethod:     tag.method,
		labelService:    tag.service,
		labelIsFailFast: strconv.FormatBool(tag.isFailFast),
	}
}