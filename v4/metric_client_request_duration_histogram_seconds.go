package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientRequestDurationHistogramVecSupportedLabels = supportedLabels{
	Code:            true,
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientRequestDurationHistogramVec instantiates client-side HistogramVec suitable for use with NewClientRequestDurationStatsHandler.
func NewClientRequestDurationHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newRequestDurationHistogramVec("client", clientRequestDurationHistogramVecSupportedLabels.labels(), opts...)
}

// ClientRequestDurationStatsHandler dedicated client-side StatsHandlerCollector that counts individual observations of request duration.
type ClientRequestDurationStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewClientRequestDurationStatsHandler instantiates ClientRequestDurationStatsHandler based on given ObserverVec and options.
// The ObserverVec must have zero, one, two, three, four or five non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service", "grpc_client_user_agent" and "grpc_code".
func NewClientRequestDurationStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientRequestDurationStatsHandler {
	h := &ClientRequestDurationStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientRequestDurationHistogramVecSupportedLabels),
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
func (h *ClientRequestDurationStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.End:
		if stat.IsClient() {
			h.vec.
				WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).
				Observe(pay.EndTime.Sub(pay.BeginTime).Seconds())
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
