package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

var serverMessageReceivedSizeHistogramVecSupportedLabels = supportedLabels{
	Method:          true,
	Service:         true,
	ClientUserAgent: true,
}

// NewServerMessageReceivedSizeHistogramVec instantiates default server-side HistogramVec suitable for use with NewServerMessageReceivedSizeStatsHandler.
func NewServerMessageReceivedSizeHistogramVec(opts ...CollectorOption) *prometheus.HistogramVec {
	return newMessageReceivedSizeHistogramVec("server", serverMessageReceivedSizeHistogramVecSupportedLabels.labels(), opts...)
}

// ServerMessageReceivedSizeStatsHandler dedicated server-side StatsHandlerCollector that counts individual observations of received message size.
type ServerMessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	serverSideLabelsHandler

	vec prometheus.ObserverVec
}

// NewServerMessageReceivedSizeStatsHandler instantiates ServerMessageReceivedSizeStatsHandler based on given ObserverVec and options.
// The ObserverVec must have zero, one, two or three non-const non-curried labels.
// For those, the only allowed names are "grpc_method", "grpc_service" and "grpc_client_user_agent".
func NewServerMessageReceivedSizeStatsHandler(vec prometheus.ObserverVec, opts ...StatsHandlerOption) *ServerMessageReceivedSizeStatsHandler {
	h := &ServerMessageReceivedSizeStatsHandler{
		vec: vec,
	}
	h.serverSideLabelsHandler = serverSideLabelsHandler{
		supportedLabels: checkLabels(vec, serverMessageReceivedSizeHistogramVecSupportedLabels),
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
func (h *ServerMessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.InPayload); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Observe(float64(pay.Length))
		}
	}
}
