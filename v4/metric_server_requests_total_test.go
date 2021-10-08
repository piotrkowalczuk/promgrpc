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

func TestNewServerRequestsTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerRequestsTotalStatsHandler(promgrpc.NewServerRequestsTotalCounterVec()))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.Begin{})
	h.HandleRPC(ctx, &stats.Begin{})
	h.HandleRPC(ctx, &stats.Begin{
		Client: true,
	})

	const metadata = `
		# HELP grpc_server_requests_received_total TODO
        # TYPE grpc_server_requests_received_total counter
	`
	expected := `
		grpc_server_requests_received_total{grpc_method="Method",grpc_service="service"} 2
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_requests_received_total"); err != nil {
		t.Fatal(err)
	}
}
