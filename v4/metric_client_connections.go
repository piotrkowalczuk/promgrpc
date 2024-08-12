package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientConnectionsGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	labels := []string{labelRemoteAddr, labelLocalAddr}
	return newConnectionsGaugeVec("client", labels, opts...)
}

type ClientConnectionsStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewClientConnectionsStatsHandler ...
func NewClientConnectionsStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ClientConnectionsStatsHandler {
	h := &ClientConnectionsStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				client: true,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)
	return h
}

// HandleConn HandleRPC processes the RPC stats.
func (h *ClientConnectionsStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.labels(ctx)...).Inc()
		}
	case *stats.ConnEnd:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.labels(ctx)...).Dec()
		}
	}
}

func (h *ClientConnectionsStatsHandler) labels(ctx context.Context) []string {
	tag := ctx.Value(tagConnKey).(connTagLabels)

	specialLabelValues := []string{
		tag.remoteAddr,
		tag.localAddr,
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
