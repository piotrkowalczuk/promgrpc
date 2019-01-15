package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
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
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
	)
}

func newRequestsGaugeVec(sub, name, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
	)
}

type RequestsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewRequestsStatsHandler ...
func NewRequestsStatsHandler(sub Subsystem, vec *prometheus.GaugeVec) *RequestsStatsHandler {
	return &RequestsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *RequestsStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC processes the RPC stats.
func (h *RequestsStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	switch stat.(type) {
	case *stats.Begin:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	case *stats.End:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Dec()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Dec()
		}
	}
}
