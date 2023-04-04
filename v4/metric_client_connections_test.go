package promgrpc_test

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/stats"
)

func TestNewClientConnectionsStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewClientConnectionsStatsHandler(promgrpc.NewClientConnectionsGaugeVec(promgrpc.CollectorWithNamespace("promgrpctest")))
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
	expected := `
		promgrpctest_client_connections{ grpc_local_addr = "1.2.3.4", grpc_remote_addr = "4.3.2.1:8080" } 2
	`
	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "promgrpctest_client_connections"); err != nil {
		t.Fatal(err)
	}
}
