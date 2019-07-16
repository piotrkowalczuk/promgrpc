package promgrpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// StatsHandlerCollector is a simple wrapper for stats Handler and prometheus Collector interfaces.
type StatsHandlerCollector interface {
	stats.Handler
	prometheus.Collector
}

var _ StatsHandlerCollector = &StatsHandler{}

// StatsHandler wraps set of stats handlers and coordinate their execution.
// Additionally, it tags RPC requests with a common set of labels.
// That way it reduces context manipulation overhead and improves overall performance.
type StatsHandler struct {
	handlers []StatsHandlerCollector
}

// NewStatsHandler allocates a new coordinator.
// It allows passing a various number of handlers that later it will iterate through.
func NewStatsHandler(handlers ...StatsHandlerCollector) *StatsHandler {
	return &StatsHandler{
		handlers: handlers,
	}
}

// ClientStatsHandler instantiates a default client-side coordinator together with every metric specific stats handler provided by this package.
func ClientStatsHandler(opts ...ShareableOption) *StatsHandler {
	return defaultStatsHandler(Client, opts...)
}

// ClientStatsHandler instantiates a default server-side coordinator together with every metric specific stats handler provided by this package.
func ServerStatsHandler(opts ...ShareableOption) *StatsHandler {
	return defaultStatsHandler(Server, opts...)
}

func defaultStatsHandler(sub Subsystem, opts ...ShareableOption) *StatsHandler {
	var (
		collectorOpts    []CollectorOption
		statsHandlerOpts []StatsHandlerOption
	)

	for _, opt := range opts {
		switch val := opt.(type) {
		case StatsHandlerOption:
			statsHandlerOpts = append(statsHandlerOpts, val)
		case CollectorOption:
			collectorOpts = append(collectorOpts, val)
		default:
			panic(fmt.Sprintf("shareable option does not implement any known type: %T", opt))
		}
	}

	return NewStatsHandler(
		NewConnectionsStatsHandler(sub, NewConnectionsGaugeVec(sub, collectorOpts...)),
		NewRequestsTotalStatsHandler(sub, NewRequestsTotalCounterVec(sub, collectorOpts...), statsHandlerOpts...),
		NewRequestsInFlightStatsHandler(sub, NewRequestsInFlightGaugeVec(sub, collectorOpts...), statsHandlerOpts...),
		NewRequestDurationStatsHandler(sub, NewRequestDurationHistogramVec(sub, collectorOpts...), statsHandlerOpts...),
		NewResponsesTotalStatsHandler(sub, NewResponsesTotalCounterVec(sub, collectorOpts...), statsHandlerOpts...),
		NewMessagesReceivedTotalStatsHandler(sub, NewMessagesReceivedTotalCounterVec(sub, collectorOpts...), statsHandlerOpts...),
		NewMessagesSentTotalStatsHandler(sub, NewMessagesSentTotalCounterVec(sub, collectorOpts...), statsHandlerOpts...),
		NewMessageSentSizeStatsHandler(sub, NewMessageSentSizeHistogramVec(sub, collectorOpts...), statsHandlerOpts...),
		NewMessageReceivedSizeStatsHandler(sub, NewMessageReceivedSizeHistogramVec(sub, collectorOpts...), statsHandlerOpts...),
	)
}

// TagRPC implements stats Handler interface.
func (h *StatsHandler) TagRPC(ctx context.Context, inf *stats.RPCTagInfo) context.Context {
	service, method := split(inf.FullMethodName)

	ctx = context.WithValue(ctx, tagRPCKey, rpcTag{
		isFailFast: strconv.FormatBool(inf.FailFast),
		service:    service,
		method:     method,
	})

	for _, c := range h.handlers {
		ctx = c.TagRPC(ctx, inf)
	}
	return ctx
}

// HandleRPC implements stats Handler interface.
func (h *StatsHandler) HandleRPC(ctx context.Context, sts stats.RPCStats) {
	for _, c := range h.handlers {
		c.HandleRPC(ctx, sts)
	}
}

// TagConn implements stats Handler interface.
func (h *StatsHandler) TagConn(ctx context.Context, inf *stats.ConnTagInfo) context.Context {
	for _, c := range h.handlers {
		ctx = c.TagConn(ctx, inf)
	}
	return ctx
}

// HandleConn implements stats Handler interface.
func (h *StatsHandler) HandleConn(ctx context.Context, sts stats.ConnStats) {
	for _, c := range h.handlers {
		c.HandleConn(ctx, sts)
	}
}

// Describe implements prometheus Collector interface.
func (h *StatsHandler) Describe(in chan<- *prometheus.Desc) {
	for _, c := range h.handlers {
		c.Describe(in)
	}
}

// Collect implements prometheus Collector interface.
func (h *StatsHandler) Collect(in chan<- prometheus.Metric) {
	for _, c := range h.handlers {
		c.Collect(in)
	}
}

type baseStatsHandler struct {
	subsystem Subsystem
	collector prometheus.Collector
	options   statsHandlerOptions
}

// TagRPC implements stats Handler interface.
func (h *baseStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

// TagConn implements stats Handler interface.
func (h *baseStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return context.WithValue(ctx, tagConnKey, prometheus.Labels{
		labelRemoteAddr:      info.RemoteAddr.String(),
		labelLocalAddr:       info.LocalAddr.String(),
		labelClientUserAgent: userAgent(ctx),
	})
}

// HandleConn implements stats Handler interface.
func (h *baseStatsHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
}

// HandleRPC implements stats Handler interface.
func (h *baseStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
}

// Describe implements prometheus Collector interface.
func (h *baseStatsHandler) Describe(in chan<- *prometheus.Desc) {
	h.collector.Describe(in)
}

// Collect implements prometheus Collector interface.
func (h *baseStatsHandler) Collect(in chan<- prometheus.Metric) {
	h.collector.Collect(in)
}
