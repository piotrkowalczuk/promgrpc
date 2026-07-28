package promgrpc

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// ShareableOption is a simple wrapper for shareable method.
// It makes it possible to distinguish options reserved for direct usage, from those that are applicable on a set of objects.
type ShareableOption interface {
	shareable()
}

type statsHandlerOptions struct {
	// stats.ConnTagInfo carries no information about whether it is an incoming or outgoing connection.
	// Use IsClient method if available.
	client           bool
	handleRPCLabelFn HandleRPCLabelFunc
	tagRPCLabelFn    TagRPCLabelFunc
}

// StatsHandlerOption configures a stats handler behaviour.
type StatsHandlerOption interface {
	apply(*statsHandlerOptions)
}

type funcStatsHandlerOption struct {
	f func(*statsHandlerOptions)
}

func (o *funcStatsHandlerOption) apply(in *statsHandlerOptions) {
	o.f(in)
}

func newFuncStatsHandlerOption(f func(*statsHandlerOptions)) *funcStatsHandlerOption {
	return &funcStatsHandlerOption{
		f: f,
	}
}

// ShareableStatsHandlerOption is StatsHandlerOption extended with shareable capability.
type ShareableStatsHandlerOption interface {
	ShareableOption
	StatsHandlerOption
}

// StatsHandlerWithHandleRPCLabelsFunc allows to inject custom HandleRPCLabelFunc to a stats handler.
// It is not shareable because there little to no chance that all stats handlers need the same set of labels.
func StatsHandlerWithHandleRPCLabelsFunc(fn HandleRPCLabelFunc) StatsHandlerOption {
	return newFuncStatsHandlerOption(func(o *statsHandlerOptions) {
		o.handleRPCLabelFn = fn
	})
}

// StatsHandlerWithTagRPCLabelsFunc allows to inject custom TagRPCLabelFunc to a stats handler.
// It is not shareable because of performance reasons.
// If all stats handlers require the same set of additional labels, it is better to implement a custom coordinator
// (e.g. by embedding StatsHandler) with self-defined TagRPC method.
// That way, it is guaranteed that new tagging execute only once and default implementation be overridden.
func StatsHandlerWithTagRPCLabelsFunc(fn TagRPCLabelFunc) StatsHandlerOption {
	return newFuncStatsHandlerOption(func(o *statsHandlerOptions) {
		o.tagRPCLabelFn = fn
	})
}

type collectorOptions struct {
	namespace           string
	userAgent           string
	constLabels         prometheus.Labels
	maxSendMsgSize      int
	maxReceivedMsgSize  int
	maxRequestDuration  float64
}

// CollectorOption configures a collector.
type CollectorOption interface {
	apply(*collectorOptions)
}

type funcCollectorOption struct {
	f func(*collectorOptions)
}

func (o *funcCollectorOption) apply(in *collectorOptions) {
	o.f(in)
}

func newFuncCollectorOption(f func(*collectorOptions)) *funcCollectorOption {
	return &funcCollectorOption{
		f: f,
	}
}

// ShareableCollectorOption is CollectorOption extended with shareable capability.
type ShareableCollectorOption interface {
	ShareableOption
	CollectorOption
}

type funcShareableCollectorOption struct {
	funcCollectorOption
}

func (o *funcShareableCollectorOption) shareable() {}

func newFuncShareableCollectorOption(f func(*collectorOptions)) *funcShareableCollectorOption {
	return &funcShareableCollectorOption{
		funcCollectorOption: funcCollectorOption{f: f},
	}
}

// CollectorWithNamespace returns a ShareableCollectorOption which sets namespace of a collector.
func CollectorWithNamespace(namespace string) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.namespace = namespace
	})
}

// CollectorWithUserAgent ...
func CollectorWithUserAgent(name, version string) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.userAgent = fmt.Sprintf("grpc/%s %s/%s", grpc.Version, name, version)
	})
}

// CollectorWithConstLabels returns a ShareableCollectorOption which adds a set of constant labels to a collector.
func CollectorWithConstLabels(constLabels prometheus.Labels) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.constLabels = constLabels
	})
}

// CollectorWithMessageReceivedMaxSize returns a ShareableCollectorOption that caps
// the generated histogram buckets for received message sizes (in bytes). If the
// size is not provided prometheus.DefBuckets is used for backward compatibility.
// Use in combination with grpc.MaxRecvMsgSize.
func CollectorWithMessageReceivedMaxSize(s int) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.maxReceivedMsgSize = s
	})
}

// CollectorWithMessageSendMaxSize returns a ShareableCollectorOption that caps the
// generated histogram buckets for sent message sizes (in bytes). If the
// size is not provided prometheus.DefBuckets is used for backward compatibility.
// Use in combination with grpc.MaxSendMsgSize.
func CollectorWithMessageSendMaxSize(s int) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.maxSendMsgSize = s
	})
}

// CollectorWithMaxRequestDuration returns a ShareableCollectorOption that caps the
// generated histogram buckets for request durations (in seconds). If the duration
// is not provided prometheus.DefBuckets is used for backward compatibility.
func CollectorWithMaxRequestDuration(d time.Duration) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.maxRequestDuration = d.Seconds()
	})
}

func applyCollectorOptions(prototype prometheus.Opts, opts ...CollectorOption) prometheus.Opts {
	var options collectorOptions
	for _, opt := range opts {
		opt.apply(&options)
	}

	if options.namespace != "" {
		prototype.Namespace = options.namespace
	}
	if options.constLabels != nil {
		prototype.ConstLabels = options.constLabels
	}

	return prototype
}

func applyHistogramOptions(prototype prometheus.HistogramOpts, opts ...CollectorOption) prometheus.HistogramOpts {
	var options collectorOptions
	for _, opt := range opts {
		opt.apply(&options)
	}

	if options.namespace != "" {
		prototype.Namespace = options.namespace
	}
	if options.constLabels != nil {
		prototype.ConstLabels = options.constLabels
	}

	return prototype
}
