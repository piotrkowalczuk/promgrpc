package promgrpc_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/alexeyxo/promgrpc/v4/pb/private/test"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func Example() {
	reg := prometheus.NewRegistry()

	ssh := promgrpc.ServerStatsHandler(
		promgrpc.CollectorWithNamespace("example"),
		promgrpc.CollectorWithConstLabels(prometheus.Labels{"service": "foo"}),
	)
	csh := promgrpc.ClientStatsHandler(
		promgrpc.CollectorWithConstLabels(prometheus.Labels{"service": "bar"}),
	)

	srv := grpc.NewServer(grpc.StatsHandler(ssh))
	imp := newDemoServer()

	test.RegisterTestServiceServer(srv, imp)
	reg.MustRegister(ssh)
	reg.MustRegister(csh)
}

func BenchmarkUnary_all(b *testing.B) {
	// Listen an actual port.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}

	reg := prometheus.NewRegistry()

	ssh := promgrpc.ServerStatsHandler()
	csh := promgrpc.ClientStatsHandler()

	srv := grpc.NewServer(grpc.StatsHandler(ssh))
	imp := newDemoServer()

	test.RegisterTestServiceServer(srv, imp)
	reg.MustRegister(ssh)
	reg.MustRegister(csh)

	go func() {
		if err := srv.Serve(lis); err != grpc.ErrServerStopped {
			b.Error(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	con, err := grpc.DialContext(ctx, lis.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStatsHandler(csh),
	)
	if err != nil {
		b.Fatal(err)
	}

	req := &test.Request{Value: "example"}
	cli := test.NewTestServiceClient(con)

	if _, err := cli.Unary(ctx, req); err != nil {
		b.Fatal(err)
	}

	ctx = context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := cli.Unary(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

// demoServiceServer defines a Server.
type demoServiceServer struct {
	test.TestServiceServer
}

func newDemoServer() *demoServiceServer {
	return &demoServiceServer{}
}

func (s *demoServiceServer) Unary(ctx context.Context, req *test.Request) (*test.Response, error) {
	return &test.Response{Value: fmt.Sprintf("unary-%s", req.GetValue())}, nil
}

func (s *demoServiceServer) ServerSide(req *test.Request, stream test.TestService_ServerSideServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&test.Response{
			Value: fmt.Sprintf("server-side-%s-%d", req.GetValue(), i),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *demoServiceServer) ClientSide(stream test.TestService_ClientSideServer) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
