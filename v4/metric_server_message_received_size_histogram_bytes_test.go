package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func TestNewServerMessageReceivedSizeStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerMessageReceivedSizeStatsHandler(promgrpc.NewServerMessageReceivedSizeHistogramVec()))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: 5,
	})

	const metadata = `
		# HELP grpc_server_message_received_size_histogram_bytes TODO
        # TYPE grpc_server_message_received_size_histogram_bytes histogram
	`
	expected := `
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="32"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="96"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="288"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="864"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2592"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="7776"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="23328"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="69984"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="209952"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="629856"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.889568e+06"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="5.668704e+06"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.7006112e+07"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="5.1018336e+07"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.53055008e+08"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="4.59165024e+08"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.377495072e+09"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="4.132485216e+09"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.2397455648e+10"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="3.7192366944e+10"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_server_message_received_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 15
        grpc_server_message_received_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_received_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
