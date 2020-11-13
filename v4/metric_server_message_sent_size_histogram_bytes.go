package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverMessageSentSizeHistogramVecSupportedLabels = supportedLabels{
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewServerMessageSentSizeHistogramVec instantiates default server-side HistogramVec suitable for use with NewServerMessageSentSizeStatsHandler.
func NewServerMessageSentSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageSentSizeHistogramVec("server", serverMessageSentSizeHistogramVecSupportedLabels.labels(), opts...)
}

// ServerMessageSentSizeStatsHandler dedicated server-side StatsHandlerCollector that counts individual observations of sent message size.
type ServerMessageSentSizeStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewServerMessageSentSizeStatsHandler instantiates ServerMessageSentSizeStatsHandler based on given ObserverVec and options.
// The ObserverVec must have zero, one, two or three non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewServerMessageSentSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerMessageSentSizeStatsHandler {
	h := &ServerMessageSentSizeStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverMessageSentSizeHistogramVecSupportedLabels),
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
func (h *ServerMessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.OutPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}
