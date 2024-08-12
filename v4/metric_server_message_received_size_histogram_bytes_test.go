package promgrpc_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"

	"github.com/piotrkowalczuk/promgrpc/v4"
)

func TestNewServerMessageReceivedSizeStatsHandler(t *testing.T) {
	ctx := promgrpc.DynamicLabelValuesToCtx(context.Background(), map[string]string{dynamicLabel: dynamicLabelValue})
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)
	h := promgrpc.NewStatsHandler(
		promgrpc.NewServerMessageReceivedSizeStatsHandler(
			promgrpc.NewServerMessageReceivedSizeHistogramVec(collectorOpts...),
			statsHandlerOpts...,
		))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: 5,
	})

	const metadata = `
		# HELP grpc_server_message_received_size_histogram_bytes TODO
        # TYPE grpc_server_message_received_size_histogram_bytes histogram
	`
	expected := fmt.Sprintf(`
		grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.005"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.01"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.025"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.05"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.1"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.25"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.5"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2.5"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="5"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="10"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_server_message_received_size_histogram_bytes_sum{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 15
        grpc_server_message_received_size_histogram_bytes_count{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`, dynamicLabel, dynamicLabelValue)

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_received_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
