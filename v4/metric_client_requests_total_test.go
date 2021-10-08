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

func TestNewClientRequestsTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientRequestsTotalStatsHandler(promgrpc.NewClientRequestsTotalCounterVec()))
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
		Client: false,
	})

	const metadata = `
		# HELP grpc_client_requests_sent_total TODO
        # TYPE grpc_client_requests_sent_total counter
	`
	expected := `
		grpc_client_requests_sent_total{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_requests_sent_total"); err != nil {
		t.Fatal(err)
	}
}
