package promgrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/piotrkowalczuk/promgrpc/v4/pb/private/test"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func suite(t *testing.T) (test.TestServiceClient, *prometheus.Registry, func(*testing.T)) {
	lis := listener(t)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ssh := promgrpc.ServerStatsHandler(
		promgrpc.CollectorWithUserAgent("test", "v4.1.2"),
		promgrpc.CollectorWithConstLabels(prometheus.Labels{"service": "test"}),
	)
	csh := promgrpc.ClientStatsHandler(
		promgrpc.CollectorWithConstLabels(prometheus.Labels{"service": "test"}),
	)
	srv := grpc.NewServer(grpc.StatsHandler(ssh))

	test.RegisterTestServiceServer(srv, newDemoServer())

	reg := prometheus.NewRegistry()
	registerCollector(t, reg, ssh)
	registerCollector(t, reg, csh)

	go func() {
		if err := srv.Serve(lis); !errors.Is(err, grpc.ErrServerStopped) {
			if err != nil {
				t.Error(err)
			}
		}
	}()

	cli, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(csh),
		grpc.WithUserAgent("test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := cli.Close(); err != nil {
			t.Error(err)
		}
	})

	return test.NewTestServiceClient(cli), reg, func(t *testing.T) {
		srv.GracefulStop()
	}
}
