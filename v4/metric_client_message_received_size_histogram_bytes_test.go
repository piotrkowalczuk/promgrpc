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

func TestNewClientMessageReceivedSizeStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(
		promgrpc.NewClientMessageReceivedSizeStatsHandler(
			promgrpc.NewClientMessageReceivedSizeHistogramVec(
				promgrpc.CollectorWithNamespace("promgrpctest"),
			),
		),
	)
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: false,
		Length: 5,
	})

	const metadata = `
		# HELP promgrpctest_client_message_received_size_histogram_bytes TODO
        # TYPE promgrpctest_client_message_received_size_histogram_bytes histogram
	`
	expected := `
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="32"} 3
		promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="59.496662720538"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="110.6204023400455"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="205.67327400112185"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="382.40229418354824"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="710.9893850187037"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1321.921738698142"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="2457.810369695954"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="4569.734831151281"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="8496.3741241032"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="15797.05955028983"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="29371.010126245237"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="54608.65885133497"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="101532.27991558747"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="188775.99416828004"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="350985.6767113877"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="652577.3883449133"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1.213318024168964e+06"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="2.2558866642728257e+06"} 3
        promgrpctest_client_message_received_size_histogram_bytes_bucket{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="4.1943039999999953e+06"} 3
		promgrpctest_client_message_received_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 15
		promgrpctest_client_message_received_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "promgrpctest_client_message_received_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
