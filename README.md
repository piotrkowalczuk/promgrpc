# promgrpc [![Build Status](https://travis-ci.org/piotrkowalczuk/promgrpc.svg?branch=master)](https://travis-ci.org/piotrkowalczuk/promgrpc)

Library allows to monitor gRPC based client and server applications.

## Metrics

### Client

* __grpc_client_errors_total__
* __grpc_client_received_messages_total__
* __grpc_client_reconnects_total__
* __grpc_client_request_duration_microseconds__
* __grpc_client_request_duration_microseconds_sum__
* __grpc_client_request_duration_microseconds_count__
* __grpc_client_requests_total__
* __grpc_client_send_messages_total__

### Server

* __grpc_server_errors_total__
* __grpc_server_received_messages_total__
* __grpc_server_request_duration_microseconds__
* __grpc_server_request_duration_microseconds_sum__
* __grpc_server_request_duration_microseconds_count__
* __grpc_server_requests_total__
* __grpc_server_send_messages_total__

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