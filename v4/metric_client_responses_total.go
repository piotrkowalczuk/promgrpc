package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/promgrpc/v4/internal/useragent"
)

// NewClientResponsesTotalCounterVec allocates a new Prometheus CounterVec for the client and given set of options.
func NewClientResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		// keep alphabetical order
		labelCode,
		labelIsFailFast,
		labelMethod,
		labelService,
		labelClientUserAgent,
	}
	return newResponsesTotalCounterVec("client", "responses_received_total", "TODO", labels, opts...)
}

// ClientResponsesTotalStatsHandler is responsible for counting number of incoming (server side) or outgoing (client side) requests.
type ClientResponsesTotalStatsHandler struct {
	baseStatsHandler
	uas useragent.Store
	vec *prometheus.CounterVec
}

// NewClientResponsesTotalStatsHandler ...
func NewClientResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientResponsesTotalStatsHandler {
	h := &ClientResponsesTotalStatsHandler{
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

// HandleRPC implements stats Handler interface.
func (h *ClientResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.End:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.OutHeader:
		_ = h.uas.ClientSide(ctx, pay)
	}
}

func (h *ClientResponsesTotalStatsHandler) labels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	specialLabelValues := []string{
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
		h.uas.ClientSide(ctx, stat),
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
