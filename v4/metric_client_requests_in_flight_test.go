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

func TestNewClientRequestsInFlightStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientRequestsInFlightStatsHandler(promgrpc.NewClientRequestsInFlightGaugeVec()))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.End{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.End{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.End{
		Client: false,
	})

	const metadata = `
		# HELP grpc_client_requests_in_flight TODO
        # TYPE grpc_client_requests_in_flight gauge
	`
	expected := `
		grpc_client_requests_in_flight{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 2
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_requests_in_flight"); err != nil {
		t.Fatal(err)
	}
}
