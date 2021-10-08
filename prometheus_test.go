package promgrpc_test

import (
	"net"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc"
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
	if err := req.Register(interceptor); err != nil {
		t.Fatal(err)
	}

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
	handler := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	ctx := handler.TagConn(context.Background(), &stats.ConnTagInfo{
		LocalAddr:  &net.TCPAddr{},
		RemoteAddr: &net.TCPAddr{},
	})

	handler.HandleConn(ctx, &stats.ConnBegin{})
	handler.HandleConn(ctx, &stats.ConnBegin{Client: true})
	handler.HandleConn(ctx, &stats.ConnEnd{})
	handler.HandleConn(ctx, &stats.ConnEnd{Client: true})
}

func TestInterceptor_HandleRPC(t *testing.T) {
	handler := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	ctx := handler.TagRPC(context.Background(), &stats.RPCTagInfo{
		FullMethodName: "method",
		FailFast:       true,
	})
	handler.HandleRPC(ctx, &stats.Begin{})
	handler.HandleRPC(ctx, &stats.Begin{Client: true})
	handler.HandleRPC(ctx, &stats.End{})
	handler.HandleRPC(ctx, &stats.End{Client: true})
}

func TestRegisterInterceptor(t *testing.T) {
	ms := mockServer{
		"test": grpc.ServiceInfo{
			Methods: []grpc.MethodInfo{
				{
					Name: "regular",
				},
				{
					Name:           "client-stream",
					IsClientStream: true,
				},
				{
					Name:           "server-stream",
					IsServerStream: true,
				},
				{
					Name:           "bidirectional-stream",
					IsClientStream: true,
					IsServerStream: true,
				},
			},
		},
	}
	interceptor1 := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
	if err := promgrpc.RegisterInterceptor(ms, interceptor1); err != nil {
		t.Fatal(err)
	}

	interceptor2 := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{TrackPeers: true})
	if err := promgrpc.RegisterInterceptor(ms, interceptor2); err != nil {
		t.Fatal(err)
	}
}

type mockServer map[string]grpc.ServiceInfo

// GetServiceInfo implements ServiceInfoProvider interface.
func (ms mockServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return ms
}
