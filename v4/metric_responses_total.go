package promgrpc

import (
	"context"
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
		[]string{
			// keep alphabetical order
			labelClientUserAgent,
			labelCode,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
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
		// labelClientUserAgent,
		// labelCode,
		// labelIsFailFast,
		// labelMethod,
		// labelService,
		lab := []string{
			tag.clientUserAgent,
			status.Code(end.Error).String(),
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
