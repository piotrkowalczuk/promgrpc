package promgrpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// NewClientResponsesTotalCounterVec allocates a new Prometheus CounterVec for the client and given set of options.
func NewClientResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		// keep alphabetical order
		labelCode,
		labelIsFailFast,
		labelMethod,
		labelService,
	}
	return newResponsesTotalCounterVec("client", "responses_received_total", "TODO", labels, opts...)
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
				handleRPCLabelFn: clientResponsesTotalLabels,
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

func clientResponsesTotalLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
