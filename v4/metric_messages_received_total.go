package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
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

// HandleRPC implements stats Handler interface.
func (h *MessagesReceivedTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.InPayload); ok {
		tag := ctx.Value(tagRPCKey).(rpcTag)
		// labelClientUserAgent,
		// labelIsFailFast,
		// labelMethod,
		// labelService,
		lab := []string{
			tag.clientUserAgent,
			tag.isFailFast,
			tag.method,
			tag.service,
		}
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(lab...).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(lab...).Inc()
		}
	}
}
