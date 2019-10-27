package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientConnectionsGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newConnectionsGaugeVec("client", opts...)
}

type ClientConnectionsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewConnectionsStatsHandler ...
func NewClientConnectionsStatsHandler(vec *prometheus.GaugeVec) *ClientConnectionsStatsHandler {
	return &ClientConnectionsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
		},
		vec: vec,
	}
}

// HandleRPC processes the RPC stats.
func (h *ClientConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case stat.IsClient():
			h.vec.With(h.labels(ctx)).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case stat.IsClient():
			h.vec.With(h.labels(ctx)).Dec()
		}
	}
}

func (h *ClientConnectionsStatsHandler) labels(ctx context.Context) prometheus.Labels {
	return ctx.Value(tagConnKey).(prometheus.Labels)
}
