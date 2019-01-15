package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func NewRequestsTotalCounterVec(sub Subsystem) *prometheus.CounterVec {
	subsystem := strings.ToLower(sub.String())
	switch sub {
	case Server:
		return newRequestsTotalCounterVec(subsystem, "requests_received_total", "TODO")
	case Client:
		return newRequestsTotalCounterVec(subsystem, "requests_sent_total", "TODO")
	default:
		// TODO: panic?
		panic("unknown subsystem")
	}
}

func newRequestsTotalCounterVec(sub, name, help string) *prometheus.CounterVec {
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

type RequestsTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed label names are "fail_fast", "handler", "service" and "user_agent".
func NewRequestsTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec) *RequestsTotalStatsHandler {
	return &RequestsTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
		},
		vec: vec,
	}
}

// Init implements StatsHandlerCollector interface.
func (h *RequestsTotalStatsHandler) Init(info map[string]grpc.ServiceInfo) error {
	return nil // TODO: implement
}

// HandleRPC implements stats Handler interface.
func (h *RequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	lab, _ := ctx.Value(tagRPCKey).(prometheus.Labels)

	if beg, ok := stat.(*stats.Begin); ok {
		switch {
		case beg.IsClient() && h.subsystem == Client:
			h.vec.With(lab).Inc()
		case !beg.IsClient() && h.subsystem == Server:
			h.vec.With(lab).Inc()
		}
	}
}
