package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func NewMessagesSentTotalCounterVec(sub Subsystem) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "messages_sent_total",
			Help:      "TODO",
		},
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
		//[]string{labelType, labelService, labelMethod}, TODO: IsServerStream and IsClientStream not available outside interceptors. Type label cannot be used.
	)
}

type MessagesSentTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewMessagesSentTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed label names are "fail_fast", "handler", "service" and "user_agent".
func NewMessagesSentTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec) *MessagesSentTotalStatsHandler {
	return &MessagesSentTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *MessagesSentTotalStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC implements stats Handler interface.
func (h *MessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	if _, ok := stat.(*stats.OutPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	}
}
