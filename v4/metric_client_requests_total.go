package promgrpc

import (
	"context"

	"github.com/piotrkowalczuk/promgrpc/v4/internal/useragent"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	labels := []string{
		labelIsFailFast,
		labelMethod,
		labelService,
		labelClientUserAgent,
	}
	return newRequestsTotalCounterVec("client", "requests_sent_total", "TODO", labels, opts...)
}

type ClientRequestsTotalStatsHandler struct {
	baseStatsHandler
	uas useragent.Store
	vec *prometheus.CounterVec
}

// NewClientRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewClientRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ClientRequestsTotalStatsHandler {
	h := &ClientRequestsTotalStatsHandler{
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
func (h *ClientRequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if !stat.IsClient() {
		return
	}
	if _, ok := stat.(*stats.OutHeader); ok {
		h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
	}
}

func (h *ClientRequestsTotalStatsHandler) labels(ctx context.Context, sts stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	specialLabelValues := []string{
		tag.isFailFast,
		tag.method,
		tag.service,
		h.uas.ClientSide(ctx, sts),
	}
	if h.options.additionalLabelValuesFn != nil {
		specialLabelValues = append(specialLabelValues, h.options.additionalLabelValuesFn(ctx)...)
	}
	return specialLabelValues
}
