package promgrpc_test

import (
	"context"
	"fmt"
	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/piotrkowalczuk/promgrpc/v4/pb/private/test"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func ExampleStatsHandler() {
	assertErr := func(err error) {
		if err != nil {
			fmt.Println("ERR:", err)
			os.Exit(1)
		}
	}
	// Listen an actual port.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assertErr(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
			assertErr(err)
		}
	}()

	con, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStatsHandler(csh))
	assertErr(err)

	for i := 0; i < 100; i++ {
		_, err = test.NewTestServiceClient(con).Unary(ctx, &test.Request{Value: "example"})
		assertErr(err)
	}

	ss, err := test.NewTestServiceClient(con).ServerSide(ctx, &test.Request{Value: "example"})
	assertErr(err)

	for {
		_, err := ss.Recv()
		if err == io.EOF {
			break
		}
		assertErr(err)
	}

	cs, err := test.NewTestServiceClient(con).ClientSide(ctx)
	assertErr(err)

	for i := 0; i < 10; i++ {
		err := cs.SendMsg(&test.Response{
			Value: fmt.Sprintf("client-side-%d", i),
		})
		assertErr(err)
	}

	srv.GracefulStop()

	mf, err := reg.Gather()
	assertErr(err)

	for _, m := range mf {
		fmt.Println(m.GetName())
	}

	// Output:
	// grpc_client_connections
	// grpc_client_message_received_size_histogram_bytes
	// grpc_client_message_sent_size_histogram_bytes
	// grpc_client_messages_received_total
	// grpc_client_messages_sent_total
	// grpc_client_request_duration_histogram_seconds
	// grpc_client_requests_in_flight
	// grpc_client_requests_sent_total
	// grpc_client_responses_received_total
	// grpc_server_connections
	// grpc_server_message_received_size_histogram_bytes
	// grpc_server_message_sent_size_histogram_bytes
	// grpc_server_messages_received_total
	// grpc_server_messages_sent_total
	// grpc_server_request_duration_histogram_seconds
	// grpc_server_requests_in_flight
	// grpc_server_requests_received_total
	// grpc_server_responses_sent_total
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

	con, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStatsHandler(csh))
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
