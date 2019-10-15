package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewRequestsInFlightGaugeVec(sub Subsystem, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub.String()),
		Name:      "requests_in_flight",
		Help:      "TODO",
	}
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			// keep alphabetical order
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

type RequestsInFlightStatsHandler struct {
	baseStatsHandler
	vec *prometheus.GaugeVec
}

// NewRequestsInFlightStatsHandler ...
func NewRequestsInFlightStatsHandler(sub Subsystem, vec *prometheus.GaugeVec, opts ...StatsHandlerOption) *RequestsInFlightStatsHandler {
	h := &RequestsInFlightStatsHandler{
		baseStatsHandler: baseStatsHandler{
			subsystem: sub,
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: requestsInFlightLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC processes the RPC stats.
func (h *RequestsInFlightStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.Begin:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	case *stats.End:
		switch {
		case stat.IsClient() && h.subsystem == Client:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Dec()
		case !stat.IsClient() && h.subsystem == Server:
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Dec()
		}
	}
}

func requestsInFlightLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	// keep alphabetical order
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
