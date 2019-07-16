package promgrpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// NewResponsesTotalCounterVec allocates a new Prometheus CounterVec for the given subsystem and set of options.
func NewResponsesTotalCounterVec(sub Subsystem, opts ...CollectorOption) *prometheus.CounterVec {
	subsystem := strings.ToLower(sub.String())
	switch sub {
	case Server:
		return newResponsesTotalCounterVec(subsystem, "responses_sent_total", "TODO", opts...)
	case Client:
		return newResponsesTotalCounterVec(subsystem, "responses_received_total", "TODO", opts...)
	default:
		// TODO: panic?
		panic("unknown subsystem")
	}
}

func newResponsesTotalCounterVec(sub, name, help string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			// keep alphabetical order
			labelClientUserAgent,
			labelCode,
			labelIsFailFast, // TODO: remove fail fast for server side
			labelMethod,
			labelService,
		},
	)
}

// ResponsesTotalStatsHandler is responsible for counting number of incoming (server side) or outgoing (client side) requests.
// 
type ResponsesTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewResponsesTotalStatsHandler ...
func NewResponsesTotalStatsHandler(sub Subsystem, vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ResponsesTotalStatsHandler {
	h := &ResponsesTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				rpcLabelFn: responsesTotalLabels,
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
func (h *ResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.End); ok {
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.rpcLabelFn(ctx, stat)...).Inc()
		}
	}
}

func responsesTotalLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
