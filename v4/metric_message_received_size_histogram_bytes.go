package promgrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewMessageReceivedSizeHistogramVec(sub Subsystem) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "message_received_size_histogram_bytes",
			Help:      "TODO",
		},
		[]string{labelIsFailFast, labelService, labelMethod, labelClientUserAgent},
	)
}

type MessageReceivedSizeStatsHandler struct {
	baseStatsHandler
	vec *prometheus.HistogramVec
}

// NewMessageReceivedSizeStatsHandler ...
func NewMessageReceivedSizeStatsHandler(sub Subsystem, vec *prometheus.HistogramVec) *MessageReceivedSizeStatsHandler {
	return &MessageReceivedSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// HandleRPC implements stats Handler interface.
func (h *MessageReceivedSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if pay, ok := stat.(*stats.InPayload); ok {
		tag := ctx.Value(tagRPCKey).(rpcTag)
		lab := prometheus.Labels{
			labelMethod:          tag.method,
			labelService:         tag.service,
			labelIsFailFast:      strconv.FormatBool(tag.isFailFast),
			labelClientUserAgent: tag.clientUserAgent,
		}
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Observe(float64(pay.Length))
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Observe(float64(pay.Length))
		}
	}
}
