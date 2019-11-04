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

	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(applyCollectorOptions(prototype, opts...)), labels,
	)
}

func newMessageReceivedSizeHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_received_size_histogram_bytes",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...), labels,
	)
}

func newMessageSentSizeHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_sent_size_histogram_bytes",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...), labels,
	)
}

func newMessagesReceivedTotalCounterVec(sub string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_received_total",
		Help:      "TODO",
	}
	return prometheus.NewCounterVec(prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)), labels)
}

func newMessagesSentTotalCounterVec(sub string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_sent_total",
		Help:      "TODO",
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		labels,
	)
}

func newRequestDurationHistogramVec(sub string, labels []string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "request_duration_histogram_seconds",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		labels,
	)
}

func newRequestsInFlightGaugeVec(sub string, labels []string, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "requests_in_flight",
		Help:      "TODO",
	}
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(applyCollectorOptions(prototype, opts...)), labels,
	)
}

func newRequestsTotalCounterVec(sub, name, help string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)), labels,
	)
}

func newResponsesTotalCounterVec(sub, name, help string, labels []string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: sub,
		Name:      name,
		Help:      help,
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		labels,
	)
}
