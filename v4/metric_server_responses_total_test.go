package promgrpc_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/promgrpc/v4"
)

func TestNewServerResponsesTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)
	h := promgrpc.NewStatsHandler(
		promgrpc.NewServerResponsesTotalStatsHandler(
			promgrpc.NewServerResponsesTotalCounterVec(collectorOpts...),
			statsHandlerOpts...,
		))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.End{
		Error: status.Error(codes.Aborted, "aborted"),
	})
	h.HandleRPC(ctx, &stats.End{})
	h.HandleRPC(ctx, &stats.End{
		Client: true,
	})

	const metadata = `
		# HELP grpc_server_responses_sent_total TODO
        # TYPE grpc_server_responses_sent_total counter
	`
	expected := fmt.Sprintf(`
		grpc_server_responses_sent_total{%[1]s="",grpc_client_user_agent="fake-user-agent",grpc_code="Aborted",grpc_method="Method",grpc_service="service"} 1
        grpc_server_responses_sent_total{%[1]s="",grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service"} 1
	`, dynamicLabel)

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_responses_sent_total"); err != nil {
		t.Fatal(err)
	}
}
