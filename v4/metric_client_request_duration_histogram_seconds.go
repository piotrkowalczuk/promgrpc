package promgrpc

import (
	"context"

	"github.com/piotrkowalczuk/promgrpc/v4/internal/useragent"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewClientRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	labels := []string{
		labelCode,
		labelIsFailFast,
		labelMethod,
		labelService,
		labelClientUserAgent,
	}
	return newRequestDurationHistogramVec("client", labels, opts...)
}

type ClientRequestDurationStatsHandler struct {
	baseStatsHandler
	uas useragent.Store
	vec prometheus.ObserverVec
}

// NewClientRequestDurationStatsHandler ...
func NewClientRequestDurationStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientRequestDurationStatsHandler {
	h := &ClientRequestDurationStatsHandler{
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
func (h *ClientRequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.End:
		if stat.IsClient() {
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(pay.EndTime.Sub(pay.BeginTime).Seconds())
		}
	case *stats.OutHeader:
		_ = h.uas.ClientSide(ctx, pay)
	}
}

func (h *ClientRequestDurationStatsHandler) labels(ctx context.Context, stat stats.RPCStats) []string {
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
