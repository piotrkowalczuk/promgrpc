package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientConnectionsGaugeVecSupportedLabels = supportedLabels{
	RemoteAddr: true,
	LocalAddr:  true,
}

// NewClientConnectionsGaugeVec instantiates client-side GaugeVec suitable for use with NewClientConnectionsStatsHandler.
func NewClientConnectionsGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newConnectionsGaugeVec("client", clientConnectionsGaugeVecSupportedLabels.labels(), opts...)
}

// ClientConnectionsStatsHandler dedicated client-side StatsHandlerCollector that counts the number of outgoing connections.
type ClientConnectionsStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.GaugeVec
}

// NewClientConnectionsStatsHandler instantiates ClientConnectionsStatsHandler based on given GaugeVec.
// The GaugeVec must have zero, one or two non-const non-curried labels.
// For those, the only allowed names are "grpc_remote_addr" and "grpc_local_addr".
func NewClientConnectionsStatsHandler(vec *prometheus.GaugeVec) *ClientConnectionsStatsHandler {
	return &ClientConnectionsStatsHandler{
		clientSideLabelsHandler: clientSideLabelsHandler{
			supportedLabels: checkLabels(vec, clientConnectionsGaugeVecSupportedLabels),
		},
		baseStatsHandler: baseStatsHandler{
			collector: vec,
		},
		vec: vec,
	}
}

// HandleConn implements stats Handler interface.
func (h *ClientConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.labelsTagConn(ctx)...).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.labelsTagConn(ctx)...).Dec()
		}
	}
}
