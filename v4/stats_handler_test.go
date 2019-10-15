package promgrpc_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc/v4/pb/private/test"
	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
)

func TestStatsHandler(t *testing.T) {
	t.Parallel()

	rpc, reg, teardown := suite(t)
	defer teardown(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type expectations struct {
		connections             float64
		clientResponsesReceived float64
		clientRequestsSent      float64
		clientMessagesSent      float64
		clientMessagesReceived  float64
		clientRequestsInFlight  float64
	}

	exp := &expectations{}

	assert := func(t *testing.T, e *expectations) {
		t.Helper()

		// CLIENT
		assertMetric(t, reg, "grpc_client_connections", e.connections)
		assertMetric(t, reg, "grpc_client_message_received_size_histogram_bytes_count", e.clientMessagesReceived)
		assertMetric(t, reg, "grpc_client_message_sent_size_histogram_bytes_count", e.clientMessagesSent)
		assertMetric(t, reg, "grpc_client_messages_received_total", e.clientMessagesReceived)
		assertMetric(t, reg, "grpc_client_messages_sent_total", e.clientMessagesSent)
		assertMetric(t, reg, "grpc_client_request_duration_histogram_seconds_count", e.clientResponsesReceived)
		assertMetric(t, reg, "grpc_client_requests_in_flight", e.clientRequestsInFlight)
		assertMetric(t, reg, "grpc_client_requests_sent_total", e.clientRequestsSent)
		assertMetric(t, reg, "grpc_client_responses_received_total", e.clientResponsesReceived)
		// SERVER
		assertMetric(t, reg, "grpc_server_connections", e.connections)
		assertMetric(t, reg, "grpc_server_message_received_size_histogram_bytes_count", e.clientMessagesSent)
		assertMetric(t, reg, "grpc_server_message_sent_size_histogram_bytes_count", e.clientMessagesReceived)
		assertMetric(t, reg, "grpc_server_messages_received_total", e.clientMessagesSent)
		assertMetric(t, reg, "grpc_server_messages_sent_total", e.clientMessagesReceived)
		assertMetric(t, reg, "grpc_server_request_duration_histogram_seconds_count", e.clientResponsesReceived)
		assertMetric(t, reg, "grpc_server_requests_in_flight", e.clientRequestsInFlight)
		assertMetric(t, reg, "grpc_server_requests_received_total", e.clientRequestsSent)
		assertMetric(t, reg, "grpc_server_responses_sent_total", e.clientResponsesReceived)
	}

	exp.connections += 1
	exp.clientRequestsSent += 100
	exp.clientResponsesReceived += 100
	exp.clientMessagesSent += 100
	exp.clientMessagesReceived += 100
	for i := 0; i < 100; i++ {
		if _, err := rpc.Unary(ctx, &test.Request{Value: "example"}); err != nil {
			t.Fatal(err)
		}
	}

	exp.clientRequestsSent += 1
	ss, err := rpc.ServerSide(ctx, &test.Request{Value: "example"})
	if err != nil {
		t.Fatal(err)
	}

	for {
		_, err := ss.Recv()
		if err == io.EOF {
			exp.clientResponsesReceived += 1
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		exp.clientMessagesReceived += 1
	}

	exp.clientRequestsSent += 1
	exp.clientMessagesSent += 1
	cs, err := rpc.ClientSide(ctx)
	if err != nil {
		t.Fatal(err)
	}

	exp.clientRequestsInFlight = 1
	t.Log("before")
	assert(t, exp)

	for i := 0; i < 10; i++ {
		err := cs.SendMsg(&test.Response{
			Value: fmt.Sprintf("client-side-%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
		exp.clientMessagesSent += 1
	}
	exp.clientResponsesReceived += 1

	teardown(t)
	exp.connections -= 1
	exp.clientRequestsInFlight = 0

	<-time.After(1 * time.Second)

	t.Log("after")
	assert(t, exp)
}

func listener(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	return lis
}

func registerCollector(t *testing.T, r *prometheus.Registry, c prometheus.Collector) {
	t.Helper()

	if err := r.Register(c); err != nil {
		t.Fatal(err)
	}
}

func assertMetric(t *testing.T, g prometheus.Gatherer, n string, exp float64) {
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
				//t.Logf("metric %s has expected value %g", n, exp)
				return
			}
		}
	}

	t.Errorf("metric %s does not exists", n)
}
