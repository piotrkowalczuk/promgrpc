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

func TestNewServerMessageSentSizeStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerMessageSentSizeStatsHandler(promgrpc.NewServerMessageSentSizeHistogramVec()))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
		Length: 5,
	})

	const metadata = `
		# HELP grpc_server_message_sent_size_histogram_bytes TODO
        # TYPE grpc_server_message_sent_size_histogram_bytes histogram
	`
	expected := `
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="32"} 3
		grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="59.496662720538"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="110.6204023400455"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="205.67327400112185"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="382.40229418354824"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="710.9893850187037"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1321.921738698142"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2457.810369695954"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="4569.734831151281"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="8496.3741241032"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="15797.05955028983"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="29371.010126245237"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="54608.65885133497"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="101532.27991558747"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="188775.99416828004"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="350985.6767113877"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="652577.3883449133"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1.213318024168964e+06"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2.2558866642728257e+06"} 3
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="4.1943039999999953e+06"} 3
		grpc_server_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
		grpc_server_message_sent_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 15
		grpc_server_message_sent_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_sent_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}

func TestNewServerMessageSentSizeStatsHandler_300KB(t *testing.T) {
	const size = 300 * 1024
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewServerMessageSentSizeStatsHandler(
		promgrpc.NewServerMessageSentSizeHistogramVec(
			promgrpc.CollectorWithMessageSendMaxSize(size),
		),
	))
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Length: size,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
		Length: size,
	})

	const metadata = `
		# HELP grpc_server_message_sent_size_histogram_bytes TODO
        # TYPE grpc_server_message_sent_size_histogram_bytes histogram
	`
	expected := `
		grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="32"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="51.84933624733461"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="84.01105216528647"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="136.1224153815721"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="220.55802768495266"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="357.368354358939"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="579.0409990410393"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="938.2153581334813"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="1520.1826117586359"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="2463.139356075776"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="3991.004396788006"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="6466.5915291766305"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="10477.764955326402"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="16977.0361965393"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="27507.75181992303"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="44570.58354748186"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="72217.34915916585"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="117013.17560764868"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="189595.48398279343"} 0
        grpc_server_message_sent_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="307200"} 3
		grpc_server_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service",le="+Inf"} 3
		grpc_server_message_sent_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 921600
		grpc_server_message_sent_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_message_sent_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
