package promgrpc

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func newConnectionsGaugeVec(sub string, labels []string, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "connections",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewGaugeVec(prometheus.GaugeOpts(promOpts), append(labels, additionalDynamicLabels...))
}

func newMessageReceivedSizeHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_received_size_histogram_bytes",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyHistogramOptions(prototype, opts...)
	return prometheus.NewHistogramVec(promOpts, append(labels, additionalDynamicLabels...))
}

func newMessageSentSizeHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_sent_size_histogram_bytes",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyHistogramOptions(prototype, opts...)
	return prometheus.NewHistogramVec(promOpts, append(labels, additionalDynamicLabels...))
}

func newMessagesReceivedTotalCounterVec(sub string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_received_total",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewCounterVec(prometheus.CounterOpts(promOpts), append(labels, additionalDynamicLabels...))
}

func newMessagesSentTotalCounterVec(sub string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_sent_total",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(promOpts),
		append(labels, additionalDynamicLabels...),
	)
}

func newRequestDurationHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "request_duration_histogram_seconds",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyHistogramOptions(prototype, opts...)
	return prometheus.NewHistogramVec(promOpts, append(labels, additionalDynamicLabels...))
}

func newRequestsInFlightGaugeVec(sub string, labels []string, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "requests_in_flight",
		Help:      "TODO",
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(promOpts), append(labels, additionalDynamicLabels...),
	)
}

func newRequestsTotalCounterVec(sub, name, help string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(promOpts), append(labels, additionalDynamicLabels...),
	)
}

func newResponsesTotalCounterVec(sub, name, help string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	promOpts, additionalDynamicLabels := applyCollectorOptions(prototype, opts...)
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(promOpts),
		append(labels, additionalDynamicLabels...),
	)
}
