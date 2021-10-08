package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

func TestNewServerResponsesTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerResponsesTotalStatsHandler(promgrpc.NewServerResponsesTotalCounterVec()))
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
	expected := `
		grpc_server_responses_sent_total{grpc_client_user_agent="fake-user-agent",grpc_code="Aborted",grpc_method="Method",grpc_service="service"} 1
        grpc_server_responses_sent_total{grpc_client_user_agent="fake-user-agent",grpc_code="OK",grpc_method="Method",grpc_service="service"} 1
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_responses_sent_total"); err != nil {
		t.Fatal(err)
	}
}
