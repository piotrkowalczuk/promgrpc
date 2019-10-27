package promgrpc

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func newConnectionsGaugeVec(sub string, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "connections",
		Help:      "TODO",
	}

	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts(applyCollectorOptions(prototype, opts...)),
		[]string{labelRemoteAddr, labelLocalAddr, labelClientUserAgent},
	)
}

func newMessageReceivedSizeHistogramVec(sub string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_received_size_histogram_bytes",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

func newMessageSentSizeHistogramVec(sub string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "message_sent_size_histogram_bytes",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

func newMessagesReceivedTotalCounterVec(sub string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_received_total",
		Help:      "TODO",
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

func newMessagesSentTotalCounterVec(sub string, opts ...CollectorOption) *prometheus.CounterVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "messages_sent_total",
		Help:      "TODO",
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts(applyCollectorOptions(prototype, opts...)),
		[]string{
			labelClientUserAgent,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

func newRequestDurationHistogramVec(sub string, opts ...CollectorOption) *prometheus.HistogramVec {
	prototype := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
		Name:      "request_duration_histogram_seconds",
		Help:      "TODO",
	}
	return prometheus.NewHistogramVec(
		applyHistogramOptions(prototype, opts...),
		[]string{
			labelClientUserAgent,
			labelCode,
			labelIsFailFast,
			labelMethod,
			labelService,
		},
	)
}

func newRequestsInFlightGaugeVec(sub string, opts ...CollectorOption) *prometheus.GaugeVec {
	prototype := prometheus.Opts{
		Namespace: namespace,
		Subsystem: strings.ToLower(sub),
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
