package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/stats"
)

func TestNewClientRequestDurationStatsHandler_30s(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bt := time.Now()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientRequestDurationStatsHandler(
		promgrpc.NewClientRequestDurationHistogramVec(
			promgrpc.CollectorWithMaxRequestDuration(30 * time.Second),
		),
	))
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
		EndTime:   bt.Add(25 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(25 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    true,
		BeginTime: bt,
		EndTime:   bt.Add(25 * time.Second),
	})
	h.HandleRPC(ctx, &stats.End{
		Client:    false,
		BeginTime: bt,
		EndTime:   bt.Add(25 * time.Second),
	})

	const md = `
		# HELP grpc_client_request_duration_histogram_seconds TODO
        # TYPE grpc_client_request_duration_histogram_seconds histogram
	`
	expected := `
		grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.001"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.001720433778486175"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.002959892386156217"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.0050922988418272"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.008760962937625542"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.015072656569956452"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.025931507494474645"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.04461344142056158"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.07675447159444837"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.1320509855809466"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.22718497607585136"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.3908567068054682"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.6724430809359947"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1.1568937905515981"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1.9903591553858793"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="3.42428112224508"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="5.891248909742982"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="10.13550362179168"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="17.437462792899368"} 0
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="30"} 3
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_client_request_duration_histogram_seconds_sum{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 75
        grpc_client_request_duration_histogram_seconds_count{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(md+expected), "grpc_client_request_duration_histogram_seconds"); err != nil {
		t.Fatal(err)
	}
}

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

// TestNewClientRequestDurationStatsHandler_1ms covers a max duration equal to the
// lowest generated bucket, which used to panic on the first observed request.
func TestNewClientRequestDurationStatsHandler_1ms(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bt := time.Now()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientRequestDurationStatsHandler(
		promgrpc.NewClientRequestDurationHistogramVec(
			promgrpc.CollectorWithMaxRequestDuration(time.Millisecond),
		),
	))
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
		EndTime:   bt.Add(time.Millisecond),
	})

	const md = `
		# HELP grpc_client_request_duration_histogram_seconds TODO
        # TYPE grpc_client_request_duration_histogram_seconds histogram
	`
	expected := `
		grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.001"} 1
        grpc_client_request_duration_histogram_seconds_bucket{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="+Inf"} 1
        grpc_client_request_duration_histogram_seconds_sum{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 0.001
        grpc_client_request_duration_histogram_seconds_count{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 1
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(md+expected), "grpc_client_request_duration_histogram_seconds"); err != nil {
		t.Fatal(err)
	}
}
