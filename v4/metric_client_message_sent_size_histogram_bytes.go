package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var clientMessageSentSizeHistogramVecSupportedLabels = supportedLabels{
	IsFailFast:      true,
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewClientMessageSentSizeHistogramVec instantiates client-side HistogramVec suitable for use with NewClientMessageSentSizeStatsHandler.
func NewClientMessageSentSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageSentSizeHistogramVec("client", clientMessageSentSizeHistogramVecSupportedLabels.labels(), opts...)
}

// ClientMessageSentSizeStatsHandler dedicated client-side StatsHandlerCollector that counts individual observations of sent message size.
type ClientMessageSentSizeStatsHandler struct {
	baseStatsHandler
	clientSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewClientMessageSentSizeStatsHandler instantiates ClientMessageSentSizeStatsHandler based on given ObserverVec and options.
// The ObserverVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed names are "grpc_is_fail_fast", "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewClientMessageSentSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ClientMessageSentSizeStatsHandler {
	h := &ClientMessageSentSizeStatsHandler{
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
func (h *ClientMessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch pay := stat.(type) {
	case *stats.OutPayload:
		if stat.IsClient() {
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	case *stats.OutHeader:
		_ = h.clientSideLabelsHandler.userAgentStore.ClientSide(ctx, pay)
	}
}
