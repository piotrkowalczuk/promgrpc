package promgrpc_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/stats"

	"github.com/piotrkowalczuk/promgrpc/v4"
)

func TestNewClientConnectionsStatsHandler(t *testing.T) {
	ctx := promgrpc.DynamicLabelValuesToCtx(context.Background(), map[string]string{dynamicLabel: dynamicLabelValue})
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorWithNamespace("promgrpctest"),
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)

	h := promgrpc.NewClientConnectionsStatsHandler(promgrpc.NewClientConnectionsGaugeVec(collectorOpts...), statsHandlerOpts...)
	ctx = h.TagConn(ctx, &stats.ConnTagInfo{
		LocalAddr: &net.TCPAddr{
			IP:   net.IPv4(1, 2, 3, 4),
			Port: 4213412,
			Zone: "",
		},
		RemoteAddr: &net.TCPAddr{
			IP:   net.IPv4(4, 3, 2, 1),
			Port: 8080,
			Zone: "",
		},
	})
	h.HandleConn(ctx, &stats.ConnBegin{
		Client: true,
	})
	ctx = h.TagConn(ctx, &stats.ConnTagInfo{
		LocalAddr: &net.TCPAddr{
			IP:   net.IPv4(1, 2, 3, 4),
			Port: 543543,
			Zone: "",
		},
		RemoteAddr: &net.TCPAddr{
			IP:   net.IPv4(4, 3, 2, 1),
			Port: 8080,
			Zone: "",
		},
	})
	h.HandleConn(ctx, &stats.ConnBegin{
		Client: true,
	})

	const metadata = `
		# HELP promgrpctest_client_connections TODO
		# TYPE promgrpctest_client_connections gauge
	`
	expected := fmt.Sprintf(`
		promgrpctest_client_connections{ %s="%s",grpc_local_addr = "1.2.3.4", grpc_remote_addr = "4.3.2.1:8080" } 2
	`, dynamicLabel, dynamicLabelValue)
	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "promgrpctest_client_connections"); err != nil {
		t.Fatal(err)
	}
}
