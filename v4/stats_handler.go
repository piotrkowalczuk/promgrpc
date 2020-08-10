package promgrpc

import (
	"context"
	"fmt"
	"net"
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
	collectorOpts, statsHandlerOpts := optionsSplit(opts...)

	return NewStatsHandler(
		NewClientConnectionsStatsHandler(NewClientConnectionsGaugeVec(collectorOpts...)),
		NewClientRequestsTotalStatsHandler(NewClientRequestsTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewClientRequestsInFlightStatsHandler(NewClientRequestsInFlightGaugeVec(collectorOpts...), statsHandlerOpts...),
		NewClientRequestDurationStatsHandler(NewClientRequestDurationHistogramVec(collectorOpts...), statsHandlerOpts...),
		NewClientResponsesTotalStatsHandler(NewClientResponsesTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewClientMessagesReceivedTotalStatsHandler(NewClientMessagesReceivedTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewClientMessagesSentTotalStatsHandler(NewClientMessagesSentTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewClientMessageSentSizeStatsHandler(NewClientMessageSentSizeHistogramVec(collectorOpts...), statsHandlerOpts...),
		NewClientMessageReceivedSizeStatsHandler(NewClientMessageReceivedSizeHistogramVec(collectorOpts...), statsHandlerOpts...),
	)
}

// ClientStatsHandler instantiates a default server-side coordinator together with every metric specific stats handler provided by this package.
func ServerStatsHandler(opts ...ShareableOption) *StatsHandler {
	collectorOpts, statsHandlerOpts := optionsSplit(opts...)

	return NewStatsHandler(
		NewServerConnectionsStatsHandler(NewServerConnectionsGaugeVec(collectorOpts...)),
		NewServerRequestsTotalStatsHandler(NewServerRequestsTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewServerRequestsInFlightStatsHandler(NewServerRequestsInFlightGaugeVec(collectorOpts...), statsHandlerOpts...),
		NewServerRequestDurationStatsHandler(NewServerRequestDurationHistogramVec(collectorOpts...), statsHandlerOpts...),
		NewServerResponsesTotalStatsHandler(NewServerResponsesTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewServerMessagesReceivedTotalStatsHandler(NewServerMessagesReceivedTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewServerMessagesSentTotalStatsHandler(NewServerMessagesSentTotalCounterVec(collectorOpts...), statsHandlerOpts...),
		NewServerMessageSentSizeStatsHandler(NewServerMessageSentSizeHistogramVec(collectorOpts...), statsHandlerOpts...),
		NewServerMessageReceivedSizeStatsHandler(NewServerMessageReceivedSizeHistogramVec(collectorOpts...), statsHandlerOpts...),
	)
}

// TagRPC implements stats Handler interface.
func (h *StatsHandler) TagRPC(ctx context.Context, inf *stats.RPCTagInfo) context.Context {
	service, method := split(inf.FullMethodName)

	ctx = context.WithValue(ctx, tagRPCKey, rpcTagLabels{
		isFailFast:      strconv.FormatBool(inf.FailFast),
		service:         service,
		method:          method,
		clientUserAgent: userAgentOnServerSide(ctx, inf),
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
	collector prometheus.Collector
	options   statsHandlerOptions
}

// TagRPC implements stats Handler interface.
func (h *baseStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if h.options.tagRPCLabelFn != nil {
		return h.options.tagRPCLabelFn(ctx, info)
	}
	return ctx
}

// TagConn implements stats Handler interface.
func (h *baseStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	remoteAddr, _, _ := net.SplitHostPort(info.RemoteAddr.String())

	return context.WithValue(ctx, tagConnKey, connTagLabels{
		remoteAddr:      remoteAddr,
		localAddr:       info.LocalAddr.String(),
		clientUserAgent: userAgentOnServerSide(ctx, &stats.RPCTagInfo{}),
	})
}

// HandleConn implements stats Handler interface.
func (h *baseStatsHandler) HandleConn(_ context.Context, _ stats.ConnStats) {
}

// HandleRPC implements stats Handler interface.
func (h *baseStatsHandler) HandleRPC(_ context.Context, _ stats.RPCStats) {
}

// Describe implements prometheus Collector interface.
func (h *baseStatsHandler) Describe(in chan<- *prometheus.Desc) {
	h.collector.Describe(in)
}

// Collect implements prometheus Collector interface.
func (h *baseStatsHandler) Collect(in chan<- prometheus.Metric) {
	h.collector.Collect(in)
}

func (h *baseStatsHandler) applyOpts(opts ...StatsHandlerOption) {
	for _, opt := range opts {
		opt.apply(&h.options)
	}
}

func optionsSplit(opts ...ShareableOption) ([]CollectorOption, []StatsHandlerOption) {
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

	return collectorOpts, statsHandlerOpts
}
