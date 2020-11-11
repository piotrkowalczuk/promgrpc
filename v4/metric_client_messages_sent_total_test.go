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

func TestNewClientMessagesSentTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	g := promgrpc.NewClientMessagesSentTotalCounterVec()
	h := promgrpc.NewStatsHandler(promgrpc.NewClientMessagesSentTotalStatsHandler(g))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: false,
	})

	<-time.After(100 * time.Millisecond)

	const metadata = `
		# HELP grpc_client_messages_sent_total TODO
        # TYPE grpc_client_messages_sent_total counter
	`
	expected := `
		grpc_client_messages_sent_total{grpc_client_user_agent="n/a/y",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(g, strings.NewReader(metadata+expected), "grpc_client_messages_sent_total"); err != nil {
		t.Fatal(err)
	}
}
