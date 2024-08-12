package promgrpc_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/stats"

	"github.com/piotrkowalczuk/promgrpc/v4"
)

func TestNewServerRequestDurationStatsHandler(t *testing.T) {
	ctx := promgrpc.DynamicLabelValuesToCtx(context.Background(), map[string]string{dynamicLabel: dynamicLabelValue})
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	bt := time.Now()
	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)
	h := promgrpc.NewStatsHandler(
		promgrpc.NewServerRequestDurationStatsHandler(
			promgrpc.NewServerRequestDurationHistogramVec(collectorOpts...),
			statsHandlerOpts...,
		))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.End{
		BeginTime: bt,
		EndTime:   bt.Add(5 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		BeginTime: bt,
		EndTime:   bt.Add(4 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		BeginTime: bt,
		EndTime:   bt.Add(3 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(1 * time.Second),
	})

	const metadata = `
		# HELP grpc_server_request_duration_histogram_seconds TODO
        # TYPE grpc_server_request_duration_histogram_seconds histogram
	`
	expected := fmt.Sprintf(`
		grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.005"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.01"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.025"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.05"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.1"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.25"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="0.5"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="1"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="2.5"} 0
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="5"} 3
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="10"} 3
        grpc_server_request_duration_histogram_seconds_bucket{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_server_request_duration_histogram_seconds_sum{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service"} 12
        grpc_server_request_duration_histogram_seconds_count{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service"} 3
	`, dynamicLabel, dynamicLabelValue)

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_request_duration_histogram_seconds"); err != nil {
		t.Fatal(err)
	}
}
