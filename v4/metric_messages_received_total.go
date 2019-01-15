package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func NewMessagesReceivedTotalCounterVec(sub Subsystem) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: strings.ToLower(sub.String()),
			Name:      "messages_received_total",
			Help:      "TODO",
		},
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
		//[]string{labelType, labelService, labelMethod}, TODO: IsServerStream and IsClientStream not available outside interceptors. Type label cannot be used.
	)
}

type MessagesReceivedTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewMessagesReceivedTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed label names are "fail_fast", "handler", "service" and "user_agent".
func NewMessagesReceivedTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec) *MessagesReceivedTotalStatsHandler {
	return &MessagesReceivedTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *MessagesReceivedTotalStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC implements stats Handler interface.
func (h *MessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	if _, ok := stat.(*stats.InPayload); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	}
}
