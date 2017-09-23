# promgrpc [![Build Status](https://travis-ci.org/piotrkowalczuk/promgrpc.svg?branch=master)](https://travis-ci.org/piotrkowalczuk/promgrpc)

[![GoDoc](https://godoc.org/github.com/piotrkowalczuk/promgrpc?status.svg)](http://godoc.org/github.com/piotrkowalczuk/promgrpc)

Library allows to monitor gRPC based client and server applications.

## Metrics

### Client

* grpc_client_connections
* grpc_client_reconnects_total
* grpc_client_errors_total
* grpc_client_requests
* grpc_client_requests_total
* grpc_client_request_duration_microseconds
* grpc_client_request_duration_microseconds_sum
* grpc_client_request_duration_microseconds_count
* grpc_client_received_messages_total
* grpc_client_send_messages_total

### Server

* grpc_server_connections
* grpc_server_errors_total
* grpc_server_requests
* grpc_server_requests_total
* grpc_server_request_duration_microseconds
* grpc_server_request_duration_microseconds_sum
* grpc_server_request_duration_microseconds_count
* grpc_server_received_messages_total
* grpc_server_send_messages_total

## Example

```go
ict := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
dop := []grpc.DialOption{
	grpc.WithDialer(inter.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", addr, timeout)
	})),
	grpc.WithStreamInterceptor(ict.StreamClient()),
	grpc.WithUnaryInterceptor(ict.UnaryClient()),
}

sop := []grpc.ServerOption{
	grpc.StatsHandler(ict),
	grpc.StreamInterceptor(ict.StreamServer()),
	grpc.UnaryInterceptor(ict.UnaryServer()),
}

prometheus.DefaultRegisterer.Register(ict)
```