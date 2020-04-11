package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
)

func AssertMetricValue(t *testing.T, g prometheus.Gatherer, n string, exp float64) {
	t.Helper()

	mf, err := g.Gather()
	if err != nil {
		t.Fatal(err)
	}

	name := strings.TrimSuffix(n, "_sum")
	name = strings.TrimSuffix(name, "_count")

	for _, m := range mf {
		if m.GetName() == name {
			var got float64
			for _, metric := range m.GetMetric() {
				switch m.GetType() {
				case prommodel.MetricType_COUNTER:
					got += metric.GetCounter().GetValue()
				case prommodel.MetricType_GAUGE:
					got += metric.GetGauge().GetValue()
				case prommodel.MetricType_HISTOGRAM:
					switch {
					case strings.HasSuffix(n, "_sum"):
						got += float64(metric.GetHistogram().GetSampleSum())
					case strings.HasSuffix(n, "_count"):
						got += float64(metric.GetHistogram().GetSampleCount())
					}
				}
			}
			if got != exp {
				t.Errorf("metric %s value do not match, expected %g but got %g", n, exp, got)
				return
			} else {
				return
			}
		}
	}

	t.Errorf("metric %s does not exists", n)
}

func AssertMetricDimensions(t *testing.T, g prometheus.Gatherer, n string, dimensions map[string]string) {
	t.Helper()

	mf, err := g.Gather()
	if err != nil {
		t.Fatal(err)
	}

	name := strings.TrimSuffix(n, "_sum")
	name = strings.TrimSuffix(name, "_count")

	var sb strings.Builder
	for _, m := range mf {
		if m.GetName() == name {
			//var got float64
			for _, metric := range m.GetMetric() {
				var b strings.Builder
				var match int
			GivenLabels:
				for _, given := range metric.GetLabel() {
					for expectedKey, expectedValue := range dimensions {
						if given.GetName() != expectedKey {
							//b.WriteString(fmt.Sprintf("PAIR MISMATCH ON KEY:\nKEY:\n	%s\n	%s\nVALUE:\n	%s\n	%s\n\n", given.GetName(), expectedKey, given.GetValue(), expectedValue))
							continue
						}
						if given.GetValue() != expectedValue {
							b.WriteString(fmt.Sprintf("PAIR MISMATCH ON VALUE:\nKEY:\n	%s\n	%s\nVALUE:\n	%s\n	%s\n\n", given.GetName(), expectedKey, given.GetValue(), expectedValue))
							continue
						}
						b.WriteString(fmt.Sprintf("PAIR MATCH:\nKEY:\n	%s\n	%s\nVALUE:\n	%s\n	%s\n\n", given.GetName(), expectedKey, given.GetValue(), expectedValue))
						match++
						continue GivenLabels
					}
				}
				if match == len(dimensions) {
					return
				}
				fmt.Printf("MATCH %d:%d\n%s\n\n%s\n\n", match, len(dimensions), b.String(), metric.String())
				sb.WriteString(fmt.Sprintf("metric checked, but does not match: %s\n	%v\n	%v\n", m.GetName(), metric.GetLabel(), dimensions))
			}
		}
	}

	t.Errorf("metric %s with dimensions %v does not exists", n, dimensions)
	t.Error(sb.String())
}
