package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewConnectionsGaugeVec(sub Subsystem, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "connections",
		Help:      "TODO",
	}

	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(applyCollectorOptions(prototype, opts...)),
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

// HandleRPC processes the RPC stats.
func (h *ConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(h.labels(ctx)).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(h.labels(ctx)).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(h.labels(ctx)).Dec()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(h.labels(ctx)).Dec()
		}
	}
}

func (h *ConnectionsStatsHandler) labels(ctx context.Context) prometheus.Labels {
	return ctx.Value(tagConnKey).(prometheus.Labels)
}
