package promgrpc

import (
	"context"

	"github.com/piotrkowalczuk/promgrpc/v4/internal/useragent"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	labels := []string{
		// keep alphabetical order
		labelIsFailFast,
		labelMethod,
		labelService,
		labelClientUserAgent,
	}
	return newRequestsInFlightGaugeVec("client", labels, opts...)
}

type ClientRequestsInFlightStatsHandler struct {
	baseStatsHandler
	uas useragent.Store
	vec *prometheus.GaugeVec
}

// NewClientRequestsInFlightStatsHandler ...
func NewClientRequestsInFlightStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ClientRequestsInFlightStatsHandler {
	h := &ClientRequestsInFlightStatsHandler{
		vec: vec,
	}
	h.baseStatsHandler = baseStatsHandler{
		collector: vec,
		options: statsHandlerOptions{
			handleRPCLabelFn: h.labels,
		},
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC processes the RPC stats.
func (h *ClientRequestsInFlightStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.OutHeader:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.End:
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Dec()
		}
	}
}

func (h *ClientRequestsInFlightStatsHandler) labels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	// keep alphabetical order
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
		h.uas.ClientSide(ctx, stat),
	}
}
