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

func TestNewClientMessagesReceivedTotalStatsHandler(t *testing.T) {
	ctx := promgrpc.DynamicLabelValuesToCtx(context.Background(), map[string]string{dynamicLabel: dynamicLabelValue})
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)
	h := promgrpc.NewStatsHandler(
		promgrpc.NewClientMessagesReceivedTotalStatsHandler(
			promgrpc.NewClientMessagesReceivedTotalCounterVec(collectorOpts...),
			statsHandlerOpts...,
		))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: false,
	})

	const metadata = `
		# HELP grpc_client_messages_received_total TODO
        # TYPE grpc_client_messages_received_total counter
	`
	expected := fmt.Sprintf(`
		grpc_client_messages_received_total{ %s="%s",grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`, dynamicLabel, dynamicLabelValue)

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_messages_received_total"); err != nil {
		t.Fatal(err)
	}
}
