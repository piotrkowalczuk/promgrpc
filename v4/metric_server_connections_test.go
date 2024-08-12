package promgrpc_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"

	"github.com/piotrkowalczuk/promgrpc/v4"
)

func TestNewServerConnectionsStatsHandler(t *testing.T) {
	ctx := promgrpc.DynamicLabelValuesToCtx(context.Background(), map[string]string{dynamicLabel: dynamicLabelValue})
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()
	collectorOpts, statsHandlerOpts := promgrpc.OptionsSplit(
		promgrpc.CollectorStatsHandlerWithDynamicLabels([]string{dynamicLabel}),
	)
	h := promgrpc.NewServerConnectionsStatsHandler(
		promgrpc.NewServerConnectionsGaugeVec(collectorOpts...),
		statsHandlerOpts...,
	)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"user-agent": []string{"fake-user-agent"}})
	ctx = h.TagConn(ctx, &stats.ConnTagInfo{
		LocalAddr: &net.TCPAddr{
			IP:   net.IPv4(1, 2, 3, 4),
			Port: 90,
			Zone: "",
		},
		RemoteAddr: &net.TCPAddr{
			IP:   net.IPv4(4, 3, 2, 1),
			Port: 111,
			Zone: "",
		},
	})
	h.HandleConn(ctx, &stats.ConnBegin{})
	ctx = h.TagConn(ctx, &stats.ConnTagInfo{
		LocalAddr: &net.TCPAddr{
			IP:   net.IPv4(1, 2, 3, 4),
			Port: 80,
			Zone: "",
		},
		RemoteAddr: &net.TCPAddr{
			IP:   net.IPv4(4, 3, 2, 1),
			Port: 111,
			Zone: "",
		},
	})
	h.HandleConn(ctx, &stats.ConnBegin{})
	h.HandleConn(ctx, &stats.ConnBegin{Client: true})

	const metadata = `
		# HELP grpc_server_connections TODO
		# TYPE grpc_server_connections gauge
	`
	expected := fmt.Sprintf(`
		grpc_server_connections{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_local_addr="1.2.3.4:80",grpc_remote_addr="4.3.2.1"} 1
        grpc_server_connections{%[1]s="%[2]s",grpc_client_user_agent="fake-user-agent",grpc_local_addr="1.2.3.4:90",grpc_remote_addr="4.3.2.1"} 1
	`, dynamicLabel, dynamicLabelValue)
	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_connections"); err != nil {
		t.Fatal(err)
	}
}
