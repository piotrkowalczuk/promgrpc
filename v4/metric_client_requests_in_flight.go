package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientRequestsInFlightGaugeVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientRequestsInFlightGaugeVec instantiates client-side GaugeVec suitable for use with NewClientRequestsInFlightStatsHandler.
func NewClientRequestsInFlightGaugeVec(opts ...CollectorOption) *prometheus.GaugeVec {
	return newRequestsInFlightGaugeVec("client", clientRequestsInFlightGaugeVecSupportedLabels.labels(), opts...)
}

// ClientRequestsInFlightStatsHandler dedicated client-side StatsHandlerCollector that counts the number of requests currently in flight.
type ClientRequestsInFlightStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec *prometheus.GaugeVec
}

// NewClientRequestsInFlightStatsHandler instantiates ClientRequestsInFlightStatsHandler based on given GaugeVec and options.
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientRequestsInFlightStatsHandler(vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *ClientRequestsInFlightStatsHandler {
	h := &ClientRequestsInFlightStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientRequestsInFlightGaugeVecSupportedLabels),
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
