package promgrpc

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewResponsesTotalCounterVec(sub Subsystem) *prometheus.CounterVec {
	subsystem := strings.ToLower(sub.String())
	switch sub {
	case Server:
		return newResponsesTotalCounterVec(subsystem, "responses_sent_total", "TODO")
	case Client:
		return newResponsesTotalCounterVec(subsystem, "responses_received_total", "TODO")
	default:
		// TODO: panic?
		panic("unknown subsystem")
	}
}

func newResponsesTotalCounterVec(sub, name, help string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		[]string{labelIsFailFast, labelService, labelMethod, labelCode, labelClientUserAgent},
	)
}

type ResponsesTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewResponsesTotalStatsHandler ...
func NewResponsesTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec) *ResponsesTotalStatsHandler {
	return &ResponsesTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// HandleRPC implements stats Handler interface.
func (h *ResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if end, ok := stat.(*stats.End); ok {
		tag := ctx.Value(tagRPCKey).(rpcTag)
		lab := prometheus.Labels{
			labelMethod:          tag.method,
			labelService:         tag.service,
			labelIsFailFast:      strconv.FormatBool(tag.isFailFast),
			labelCode:            status.Code(end.Error).String(),
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
