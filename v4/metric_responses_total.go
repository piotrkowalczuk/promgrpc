package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
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
		[]string{labelFailFast, labelService, labelMethod, labelClientUserAgent},
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

// Init implements StatsHandlerCollector interface.
func (h *ResponsesTotalStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC implements stats Handler interface.
func (h *ResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	if _, ok := stat.(*stats.End); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	}
}
