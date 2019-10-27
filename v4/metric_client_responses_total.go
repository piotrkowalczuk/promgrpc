package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// NewClientResponsesTotalCounterVec allocates a new Prometheus CounterVec for the client and given set of options.
func NewClientResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newResponsesTotalCounterVec("client", "responses_received_total", "TODO", opts...)
}

// ClientResponsesTotalStatsHandler is responsible for counting number of incoming (server side) or outgoing (client side) requests.
type ClientResponsesTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewClientResponsesTotalStatsHandler ...
func NewClientResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientResponsesTotalStatsHandler {
	h := &ClientResponsesTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: responsesTotalLabels,
			},
		},
		vec: vec,
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ClientResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.End); ok {
		switch {
		case stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}
