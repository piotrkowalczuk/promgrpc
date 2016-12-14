# promgrpc [![Build Status](https://travis-ci.org/piotrkowalczuk/promgrpc.svg?branch=master)](https://travis-ci.org/piotrkowalczuk/promgrpc)

Library allows to track gRPC based client and server aplications.

## Example

```go
inter := promgrpc.NewInterceptor()
dop := []grpc.DialOption{
	grpc.WithDialer(inter.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", addr, timeout)
	})),
	grpc.WithStreamInterceptor(inter.StreamClient()),
	grpc.WithUnaryInterceptor(inter.UnaryClient()),
}

sop := []grpc.ServerOption{
	grpc.StreamInterceptor(inter.StreamServer()),
	grpc.UnaryInterceptor(inter.UnaryServer()),
}
```