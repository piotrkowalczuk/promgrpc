package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientMessageReceivedSizeHistogramVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientMessageReceivedSizeHistogramVec instantiates default client-side HistogramVec suitable for use with NewClientMessageReceivedSizeStatsHandler.
func NewClientMessageReceivedSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageReceivedSizeHistogramVec("client", clientMessageReceivedSizeHistogramVecSupportedLabels.labels(), opts...)
}

// ClientMessageReceivedSizeStatsHandler dedicated client-side StatsHandlerCollector that counts individual observations of received message size.
type ClientMessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewClientMessageReceivedSizeStatsHandler instantiates ClientMessageReceivedSizeStatsHandler based on given ObserverVec and options.
// The CounterVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientMessageReceivedSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientMessageReceivedSizeStatsHandler {
	h := &ClientMessageReceivedSizeStatsHandler{
		vec: vec,
	}
	h.clientSideLabelsHandler = clientSideLabelsHandler{
		supportedLabels: checkLabels(vec, clientMessageReceivedSizeHistogramVecSupportedLabels),
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

// HandleRPC implements stats Handler interface.
func (h *ClientMessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.InPayload:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
