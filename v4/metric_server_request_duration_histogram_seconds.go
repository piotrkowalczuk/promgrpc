package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverRequestDurationHistogramVecSupportedLabels = supportedLabels{
	Code:            true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewServerRequestDurationHistogramVec instantiates default server-side HistogramVec suitable for use with NewServerRequestDurationStatsHandler.
func NewServerRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newRequestDurationHistogramVec("server", serverRequestDurationHistogramVecSupportedLabels.labels(), opts...)
}

// ServerRequestDurationStatsHandler dedicated server-side StatsHandlerCollector that counts individual observations of request duration.
type ServerRequestDurationStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewServerRequestDurationStatsHandler instantiates ServerRequestDurationStatsHandler based on given ObserverVec and options.
// The ObserverVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service", "grpc_client_user_agent" and "grpc_code".
func NewServerRequestDurationStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerRequestDurationStatsHandler {
	h := &ServerRequestDurationStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverRequestDurationHistogramVecSupportedLabels),
	}
	h.baseStatsHandler = baseStatsHandler{
		collector: vec,
		options: statsHandlerOptions{
			handleRPCLabelFn: h.labelsTagRPC,
		},
	}
	h.applyOpts(opts...)

	return h
}

// HandleRPC processes the RPC stats.
func (h *ServerRequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		switch {
		case !stat.IsClient():
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(end.EndTime.Sub(end.BeginTime).Seconds())
		}
	}
}
