package promgrpc

import "github.com/prometheus/client_golang/prometheus"

// ShareableOption is a simple wrapper for shareable method.
// It makes it possible to distinguish options reserved for direct usage, from those that are applicable on a set of objects.
type ShareableOption interface {
	shareable()
}

type statsHandlerOptions struct {
	rpcLabelFn HandleRPCLabelFunc
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
func StatsHandlerWithHandleRPCLabelsFunc(fn HandleRPCLabelFunc) StatsHandlerOption {
	return newFuncStatsHandlerOption(func(o *statsHandlerOptions) {
		o.rpcLabelFn = fn
	})
}

type collectorOptions struct {
	namespace   string
	constLabels prometheus.Labels
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

// CollectorWithConstLabels returns a ShareableCollectorOption which adds a set of constant labels to a collector.
func CollectorWithConstLabels(constLabels prometheus.Labels) ShareableCollectorOption {
	return newFuncShareableCollectorOption(func(o *collectorOptions) {
		o.constLabels = constLabels
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
