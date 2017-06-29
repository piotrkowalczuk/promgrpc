package promgrpc

import (
	"net"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func ExampleInterceptor_Dialer() {
	interceptor := NewInterceptor(InterceptorOpts{})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDialer(interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", addr, timeout)
	})))
}

func TestInterceptor_Dialer(t *testing.T) {
	interceptor := NewInterceptor(InterceptorOpts{})
	fn := interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return nil, nil
	})
	_, err := fn("X", 1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_UnaryServer(t *testing.T) {
	interceptor := NewInterceptor(InterceptorOpts{TrackPeers: true})
	_, err := interceptor.UnaryServer()(context.Background(), nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_StreamServer(t *testing.T) {
	interceptor := NewInterceptor(InterceptorOpts{TrackPeers: true})
	err := interceptor.StreamServer()(context.Background(), nil, &grpc.StreamServerInfo{}, func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_UnaryClient(t *testing.T) {
	interceptor := NewInterceptor(InterceptorOpts{})
	err := interceptor.UnaryClient()(context.Background(), "method", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInterceptor_StreamClient(t *testing.T) {
	interceptor := NewInterceptor(InterceptorOpts{})
	_, err := interceptor.StreamClient()(context.Background(), &grpc.StreamDesc{}, nil, "method", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestRegisterInterceptor(t *testing.T) {
	interceptor1 := NewInterceptor(InterceptorOpts{})
	RegisterInterceptor(&grpc.Server{}, interceptor1)

	interceptor2 := NewInterceptor(InterceptorOpts{TrackPeers: true})
	RegisterInterceptor(&grpc.Server{}, interceptor2)
}
