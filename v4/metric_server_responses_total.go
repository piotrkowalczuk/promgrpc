package promgrpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// NewServerResponsesTotalCounterVec allocates a new Prometheus CounterVec for the server and given set of options.
func NewServerResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		// keep alphabetical order
		labelClientUserAgent,
		labelCode,
		labelMethod,
		labelService,
	}
	return newResponsesTotalCounterVec("server", "responses_sent_total", "TODO", labels, opts...)
}

// ServerResponsesTotalStatsHandler is responsible for counting number of incoming (server side) or outgoing (client side) requests.
type ServerResponsesTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerResponsesTotalStatsHandler ...
func NewServerResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerResponsesTotalStatsHandler {
	h := &ServerResponsesTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options:   statsHandlerOptions{},
		},
		vec: vec,
	}
	h.baseStatsHandler.options.handleRPCLabelFn = h.serverResponsesTotalLabels
	h.applyOpts(opts...)

	return h
}

// HandleRPC implements stats Handler interface.
func (h *ServerResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.End); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func (h *ServerResponsesTotalStatsHandler) serverResponsesTotalLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	specialLabelValues := []string{
		tag.clientUserAgent,
		status.Code(stat.(*stats.End).Error).String(),
		tag.method,
		tag.service,
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
