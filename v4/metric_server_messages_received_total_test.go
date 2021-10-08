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

func TestNewServerMessagesReceivedTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerMessagesReceivedTotalStatsHandler(promgrpc.NewServerMessagesReceivedTotalCounterVec()))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.InPayload{})
	h.HandleRPC(ctx, &stats.InPayload{})
	h.HandleRPC(ctx, &stats.InPayload{})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})

	const metadata = `
		# HELP grpc_server_messages_received_total TODO
        # TYPE grpc_server_messages_received_total counter
	`
	expected := `
		grpc_server_messages_received_total{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_messages_received_total"); err != nil {
		t.Fatal(err)
	}
}
