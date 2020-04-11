package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerConnectionsGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	labels := []string{labelRemoteAddr, labelLocalAddr, labelClientUserAgent}
	return newConnectionsGaugeVec("server", labels, opts...)
}

type ServerConnectionsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewConnectionsStatsHandler ...
func NewServerConnectionsStatsHandler(vec *prometheus.GaugeVec) *ServerConnectionsStatsHandler {
	return &ServerConnectionsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
		},
		vec: vec,
	}
}

// HandleRPC processes the RPC stats.
func (h *ServerConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.labels(ctx)...).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.labels(ctx)...).Dec()
		}
	}
}

func (h *ServerConnectionsStatsHandler) labels(ctx context.Context) []string {
	tag := ctx.Value(tagConnKey).(connTagLabels)
	return []string{
		tag.remoteAddr,
		tag.localAddr,
		tag.clientUserAgent,
	}

}
