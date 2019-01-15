package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func NewMessageSentSizeHistogramVec(sub Subsystem) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "message_sent_size_histogram_bytes",
			Help:      "TODO",
		},
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
	)
}

type MessageSentSizeStatsHandler struct {
	baseStatsHandler
	vec *prometheus.HistogramVec
}

// NewMessageSentSizeStatsHandler ...
func NewMessageSentSizeStatsHandler(sub Subsystem, vec *prometheus.HistogramVec) *MessageSentSizeStatsHandler {
	return &MessageSentSizeStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *MessageSentSizeStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC implements stats Handler interface.
func (h *MessageSentSizeStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	if pay, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Observe(float64(pay.Length))
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Observe(float64(pay.Length))
		}
	}
}
