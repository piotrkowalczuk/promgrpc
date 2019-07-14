package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewRequestsTotalCounterVec(sub Subsystem, opts ...CollectorOption) *prometheus.CounterVec {
	subsystem := strings.ToLower(sub.String())
	switch sub {
	case Server:
		return newRequestsTotalCounterVec(subsystem, "requests_received_total", "TODO", opts...)
	case Client:
		return newRequestsTotalCounterVec(subsystem, "requests_sent_total", "TODO", opts...)
	default:
		// TODO: panic?
		panic("unknown subsystem")
	}
}

func newRequestsTotalCounterVec(sub, name, help string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

// RequestsTotalLabels ...
func RequestsTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}

type RequestsTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewRequestsTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec, opts ...StatsHandlerOption) *RequestsTotalStatsHandler {
	h := &RequestsTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				rpcLabelFn: RequestsTotalLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC implements stats Handler interface.
func (h *RequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if beg, ok := stat.(*stats.Begin); ok {
		switch {
		case beg.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		case !beg.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		}
	}
}
