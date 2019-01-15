package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func NewConnectionsGaugeVec(sub Subsystem) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "connections",
			Help:      "TODO",
		},
		[]string{labelRemoteAddr, labelLocalAddr, labelClientUserAgent},
	)
}

type ConnectionsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewConnectionsStatsHandler ...
func NewConnectionsStatsHandler(sub Subsystem, vec *prometheus.GaugeVec) *ConnectionsStatsHandler {
	return &ConnectionsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *ConnectionsStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC processes the RPC stats.
func (h *ConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	lab, _ := ctx.Value(tagConnKey).(prometheus.Labels)

	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Dec()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Dec()
		}
	}
}
