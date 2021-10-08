package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func TestNewClientRequestDurationStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bt := time.Now()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientRequestDurationStatsHandler(promgrpc.NewClientRequestDurationHistogramVec()))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(5 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(4 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(3 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    false,
		BeginTime: bt,
		EndTime:   bt.Add(1 * time.Second),
	})

	const metadata = `
		# HELP grpc_client_request_duration_histogram_seconds TODO
        # TYPE grpc_client_request_duration_histogram_seconds histogram
	`
	expected := `
		grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.005"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.01"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.025"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.05"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.1"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.25"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.5"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="2.5"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="5"} 3
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="10"} 3
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_client_request_duration_histogram_seconds_sum{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 12
        grpc_client_request_duration_histogram_seconds_count{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_request_duration_histogram_seconds"); err != nil {
		t.Fatal(err)
	}
}
