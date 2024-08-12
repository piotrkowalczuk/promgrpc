package promgrpc

import (
	"context"
	"fmt"

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
	client                  bool
	handleRPCLabelFn        HandleRPCLabelFunc
	tagRPCLabelFn           TagRPCLabelFunc
	additionalLabelValuesFn AdditionalLabelValuesFunc
}

// StatsHandlerOption configures a stats handler behaviour.
type StatsHandlerOption interface {
	applyStatsHandlerOption(*statsHandlerOptions)
}

type funcStatsHandlerOption struct {
	f func(*statsHandlerOptions)
}

func (o *funcStatsHandlerOption) applyStatsHandlerOption(in *statsHandlerOptions) {
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
	namespace     string
	userAgent     string
	constLabels   prometheus.Labels
	dynamicLabels []string
}

// CollectorOption configures a collector.
type CollectorOption interface {
	applyCollectorOption(*collectorOptions)
}

type funcCollectorOption struct {
	f func(*collectorOptions)
}

func (o *funcCollectorOption) applyCollectorOption(in *collectorOptions) {
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

// ShareableCollectorStatsHandlerOption is CollectorOption and StatsHandlerOption extended with shareable capability.
type ShareableCollectorStatsHandlerOption interface {
	ShareableOption
	CollectorOption
	StatsHandlerOption
}

type funcShareableCollectorStatsHandlerOption struct {
	funcCollectorOption
	funcStatsHandlerOption
}

func (o *funcShareableCollectorStatsHandlerOption) shareable() {}

func newFuncShareableCollectorStatsHandlerOption(
	fC func(*collectorOptions),
	fSH func(options *statsHandlerOptions),
) *funcShareableCollectorStatsHandlerOption {
	return &funcShareableCollectorStatsHandlerOption{
		funcCollectorOption:    funcCollectorOption{f: fC},
		funcStatsHandlerOption: funcStatsHandlerOption{f: fSH},
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

// CollectorStatsHandlerWithDynamicLabels returns a ShareableCollectorStatsHandlerOption
// which adds a set of dynamic labels to a collector,
// provide a func to fetch label values from context
func CollectorStatsHandlerWithDynamicLabels(dynamicLabels []string) ShareableCollectorStatsHandlerOption {
	return newFuncShareableCollectorStatsHandlerOption(
		func(o *collectorOptions) {
			o.dynamicLabels = dynamicLabels
		},
		func(o *statsHandlerOptions) {
			o.additionalLabelValuesFn = func(ctx context.Context) []string {
				res := make([]string, 0, len(dynamicLabels))
				valuesFromCtx := DynamicLabelValuesFromCtx(ctx)
				for _, label := range dynamicLabels {
					res = append(res, valuesFromCtx[label])
				}
				return res
			}
		},
	)
}

func applyCollectorOptions(prototype prometheus.Opts, opts ...CollectorOption) (prometheus.Opts, []string) {
	var options collectorOptions
	for _, opt := range opts {
		opt.applyCollectorOption(&options)
	}

	if options.namespace != "" {
		prototype.Namespace = options.namespace
	}
	if options.constLabels != nil {
		prototype.ConstLabels = options.constLabels
	}

	return prototype, options.dynamicLabels
}

func applyHistogramOptions(prototype prometheus.HistogramOpts, opts ...CollectorOption) (prometheus.HistogramOpts, []string) {
	var options collectorOptions
	for _, opt := range opts {
		opt.applyCollectorOption(&options)
	}

	if options.namespace != "" {
		prototype.Namespace = options.namespace
	}
	if options.constLabels != nil {
		prototype.ConstLabels = options.constLabels
	}

	return prototype, options.dynamicLabels
}
