package promgrpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

const (
	namespace    = "grpc"
	notAvailable = "n/a"
	// magicString is used for the hacky label test in checkLabels. Remove once fixed.
	magicString = "zZgWfBxLqvG8kc8IMv3POi2Bb0tZI3vAnBx+gBaFi9FyPzB/CzKUer1yufDa"
)

type ctxKey int

const (
	tagRPCKey  ctxKey = 1
	tagConnKey ctxKey = 3
)

func split(name string) (string, string) {
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[1:i], name[i+1:]
	}
	return "unknown", "unknown"
}

func userAgentOnServerSide(ctx context.Context, _ *stats.RPCTagInfo) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua, ok := md["user-agent"]; ok && len(ua) == 1 {
			return ua[0]
		}
	}
	return notAvailable
}

func checkLabels(c prometheus.Collector, expected supportedLabels) supportedLabels {
	// TODO(piotrkowalczuk): Remove this hacky way to check for instance labels
	// 	once Descriptors can have their dimensionality queried.
	var (
		desc *prometheus.Desc
		m    prometheus.Metric
		pm   dto.Metric
		lvs  []string
	)

	// Get the Desc from the Collector.
	descc := make(chan *prometheus.Desc, 1)
	c.Describe(descc)

	select {
	case desc = <-descc:
	default:
		panic("no description provided by collector")
	}
	select {
	case <-descc:
		panic("more than one description provided by collector")
	default:
	}

	close(descc)

	// Create a ConstMetric with the Desc. Since we don't know how many
	// variable labels there are, try for as long as it needs.
	for err := errors.New("dummy"); err != nil; lvs = append(lvs, magicString) {
		m, err = prometheus.NewConstMetric(desc, prometheus.UntypedValue, 0, lvs...)
	}

	// Write out the metric into a proto message and look at the labels.
	// If the value is not the magicString, it is a constLabel, which doesn't interest us.
	// If the label is curried, it doesn't interest us.
	// In all other cases, only "code" or "method" is allowed.
	if err := m.Write(&pm); err != nil {
		panic("error checking metric for labels")
	}

	var checked supportedLabels
	for _, label := range pm.Label {
		name, value := label.GetName(), label.GetValue()
		if value != magicString || isLabelCurried(c, name) {
			continue
		}

		if !expected.isKnown(name) {
			panic("metric partitioned with non-supported labels")
		}

		checked.enable(name)
	}

	return checked
}

func isLabelCurried(c prometheus.Collector, label string) bool {
	// This is even hackier than the label test above.
	// We essentially try to curry again and see if it works.
	// But for that, we need to type-convert to the two
	// types we use here, ObserverVec or *CounterVec.
	switch v := c.(type) {
	case *prometheus.CounterVec:
		if _, err := v.CurryWith(prometheus.Labels{label: "dummy"}); err == nil {
			return false
		}
	case prometheus.ObserverVec:
		if _, err := v.CurryWith(prometheus.Labels{label: "dummy"}); err == nil {
			return false
		}
	case *prometheus.HistogramVec:
		if _, err := v.CurryWith(prometheus.Labels{label: "dummy"}); err == nil {
			return false
		}
	case *prometheus.GaugeVec:
		if _, err := v.CurryWith(prometheus.Labels{label: "dummy"}); err == nil {
			return false
		}
	default:
		panic(fmt.Sprintf("unsupported metric vec type: %T", v))
	}
	return true
}
