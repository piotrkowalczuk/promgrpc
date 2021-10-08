package promgrpc_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4/internal/testutil"
	"github.com/alexeyxo/promgrpc/v4/pb/private/test"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
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
		testutil.AssertMetricValue(t, reg, "grpc_client_connections", e.connections)
		testutil.AssertMetricValue(t, reg, "grpc_client_message_received_size_histogram_bytes_count", e.clientMessagesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_client_message_sent_size_histogram_bytes_count", e.clientMessagesSent)
		testutil.AssertMetricValue(t, reg, "grpc_client_messages_received_total", e.clientMessagesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_client_messages_sent_total", e.clientMessagesSent)
		testutil.AssertMetricValue(t, reg, "grpc_client_request_duration_histogram_seconds_count", e.clientResponsesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_client_requests_in_flight", e.clientRequestsInFlight)
		testutil.AssertMetricDimensions(t, reg, "grpc_client_requests_in_flight", map[string]string{
			"grpc_service":           "piotrkowalczuk.promgrpc.v4.test.TestService",
			"service":                "test",
			"grpc_is_fail_fast":      "true",
			"grpc_client_user_agent": fmt.Sprintf("test grpc-go/%s", grpc.Version),
		})
		testutil.AssertMetricValue(t, reg, "grpc_client_requests_sent_total", e.clientRequestsSent)
		testutil.AssertMetricDimensions(t, reg, "grpc_client_requests_sent_total", map[string]string{
			"grpc_service":           "piotrkowalczuk.promgrpc.v4.test.TestService",
			"service":                "test",
			"grpc_is_fail_fast":      "true",
			"grpc_client_user_agent": fmt.Sprintf("test grpc-go/%s", grpc.Version),
		})
		testutil.AssertMetricValue(t, reg, "grpc_client_responses_received_total", e.clientResponsesReceived)
		// SERVER
		testutil.AssertMetricValue(t, reg, "grpc_server_connections", e.connections)
		testutil.AssertMetricValue(t, reg, "grpc_server_message_received_size_histogram_bytes_count", e.clientMessagesSent)
		testutil.AssertMetricDimensions(t, reg, "grpc_server_message_received_size_histogram_bytes_count", map[string]string{
			"grpc_service":           "piotrkowalczuk.promgrpc.v4.test.TestService",
			"service":                "test",
			"grpc_client_user_agent": fmt.Sprintf("test grpc-go/%s", grpc.Version),
		})
		testutil.AssertMetricValue(t, reg, "grpc_server_message_sent_size_histogram_bytes_count", e.clientMessagesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_server_messages_received_total", e.clientMessagesSent)
		testutil.AssertMetricValue(t, reg, "grpc_server_messages_sent_total", e.clientMessagesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_server_request_duration_histogram_seconds_count", e.clientResponsesReceived)
		testutil.AssertMetricValue(t, reg, "grpc_server_requests_in_flight", e.clientRequestsInFlight)
		testutil.AssertMetricValue(t, reg, "grpc_server_requests_received_total", e.clientRequestsSent)
		testutil.AssertMetricValue(t, reg, "grpc_server_responses_sent_total", e.clientResponsesReceived)
		testutil.AssertMetricDimensions(t, reg, "grpc_server_message_received_size_histogram_bytes_count", map[string]string{
			"grpc_service":           "piotrkowalczuk.promgrpc.v4.test.TestService",
			"service":                "test",
			"grpc_client_user_agent": fmt.Sprintf("test grpc-go/%s", grpc.Version),
		})

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
