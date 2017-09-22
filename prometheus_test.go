package promgrpc_test

import (
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/promgrpc"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func ExampleInterceptor_Dialer() {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDialer(interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", addr, timeout)
	})))
}

func TestIntercepror_Collector(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	req := prometheus.NewRegistry()
	req.Register(interceptor)

	_, err := req.Gather()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_Dialer(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	fn := interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return nil, nil
	})
	_, err := fn("X", 1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_UnaryServer(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{TrackPeers: true})
	_, err := interceptor.UnaryServer()(context.Background(), nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_StreamServer(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{TrackPeers: true})
	err := interceptor.StreamServer()(context.Background(), nil, &grpc.StreamServerInfo{}, func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_UnaryClient(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	err := interceptor.UnaryClient()(context.Background(), "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_StreamClient(t *testing.T) {
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	_, err := interceptor.StreamClient()(context.Background(), &grpc.StreamDesc{}, nil, "method", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_HandleConn(t *testing.T) {
	var handler stats.Handler
	handler = promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	handler.HandleConn(context.Background(), &stats.ConnBegin{})
	handler.HandleConn(context.Background(), &stats.ConnBegin{Client: true})
	handler.HandleConn(context.Background(), &stats.ConnEnd{})
	handler.HandleConn(context.Background(), &stats.ConnEnd{Client: true})
}

func TestRegisterInterceptor(t *testing.T) {
	interceptor1 := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	promgrpc.RegisterInterceptor(&grpc.Server{}, interceptor1)

	interceptor2 := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{TrackPeers: true})
	promgrpc.RegisterInterceptor(&grpc.Server{}, interceptor2)
}
