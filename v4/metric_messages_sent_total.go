package promgrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
		[]string{labelIsFailFast, labelService, labelMethod, labelClientUserAgent},
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

// HandleRPC implements stats Handler interface.
func (h *MessagesSentTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.OutPayload); ok {
		tag := ctx.Value(tagRPCKey).(rpcTag)
		lab := prometheus.Labels{
			labelMethod:          tag.method,
			labelService:         tag.service,
			labelIsFailFast:      strconv.FormatBool(tag.isFailFast),
			labelClientUserAgent: tag.clientUserAgent,
		}
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	}
}
