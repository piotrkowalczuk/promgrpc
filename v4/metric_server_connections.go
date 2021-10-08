package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverConnectionsGaugeVecSupportedLabels = supportedLabels{
	RemoteAddr:      true,
	LocalAddr:       true,
	ClientUserAgent: true,
}

// NewServerConnectionsGaugeVec instantiates client-side GaugeVec suitable for use with NewServerConnectionsStatsHandler.
func NewServerConnectionsGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newConnectionsGaugeVec("server", serverConnectionsGaugeVecSupportedLabels.labels(), opts...)
}

// ServerConnectionsStatsHandler dedicated server-side StatsHandlerCollector that counts the number of incoming connections.
type ServerConnectionsStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec *prometheus.GaugeVec
}

// NewServerConnectionsStatsHandler instantiates ServerConnectionsStatsHandler based on given GaugeVec.
// The GaugeVec must have zero, one, two or three non-const non-curried labels.
// For those, the only allowed names are "grpc_remote_addr", "grpc_local_addr" and "grpc_client_user_agent".
func NewServerConnectionsStatsHandler(vec *prometheus.GaugeVec) *ServerConnectionsStatsHandler {
	return &ServerConnectionsStatsHandler{
		serverSideLabelsHandler: serverSideLabelsHandler{
			supportedLabels: checkLabels(vec, serverConnectionsGaugeVecSupportedLabels),
		},
		baseStatsHandler: baseStatsHandler{
			collector: vec,
		},
		vec: vec,
	}
}

// HandleConn implements stats Handler interface.
func (h *ServerConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.labelsTagConn(ctx)...).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.labelsTagConn(ctx)...).Dec()
		}
	}
}
