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
		grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.005"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.01"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.025"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.05"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.1"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.25"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="0.5"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2.5"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="5"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="10"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_server_message_received_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 15
        grpc_server_message_received_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_received_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}

func TestNewServerMessageReceivedSizeStatsHandler_300KB(t *testing.T) {
	const size = 300 * 1024
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerMessageReceivedSizeStatsHandler(
		promgrpc.NewServerMessageReceivedSizeHistogramVec(
			promgrpc.CollectorWithMessageReceivedMaxSize(size),
		),
	))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: size,
	})

	const metadata = `
		# HELP grpc_server_message_received_size_histogram_bytes TODO
        # TYPE grpc_server_message_received_size_histogram_bytes histogram
	`
	expected := `
		grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="32"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="51.84933624733461"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="84.01105216528647"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="136.1224153815721"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="220.55802768495266"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="357.368354358939"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="579.0409990410393"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="938.2153581334813"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1520.1826117586359"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2463.139356075776"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="3991.004396788006"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="6466.5915291766305"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="10477.764955326402"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="16977.0361965393"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="27507.75181992303"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="44570.58354748186"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="72217.34915916585"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="117013.17560764868"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="189595.48398279343"} 0
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="307200"} 3
        grpc_server_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
        grpc_server_message_received_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 921600
        grpc_server_message_received_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_received_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
