package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/stats"
)

func TestNewClientMessagesReceivedTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientMessagesReceivedTotalStatsHandler(promgrpc.NewClientMessagesReceivedTotalCounterVec()))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
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

	<-time.After(100 * time.Millisecond)

	const metadata = `
		# HELP grpc_client_messages_received_total TODO
        # TYPE grpc_client_messages_received_total counter
	`
	expected := `
		grpc_client_messages_received_total{grpc_client_user_agent="n/a/y",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_messages_received_total"); err != nil {
		t.Fatal(err)
	}
}
