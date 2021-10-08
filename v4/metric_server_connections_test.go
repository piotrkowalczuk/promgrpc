package promgrpc_test

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func TestNewServerConnectionsStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewServerConnectionsStatsHandler(promgrpc.NewServerConnectionsGaugeVec())
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
	expected := `
		grpc_server_connections{grpc_client_user_agent="fake-user-agent",grpc_local_addr="1.2.3.4:80",grpc_remote_addr="4.3.2.1"} 1
        grpc_server_connections{grpc_client_user_agent="fake-user-agent",grpc_local_addr="1.2.3.4:90",grpc_remote_addr="4.3.2.1"} 1
	`
	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_server_connections"); err != nil {
		t.Fatal(err)
	}
}
