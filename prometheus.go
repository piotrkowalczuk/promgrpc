package promgrpc

import (
	"net"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func RegisterInterceptor(s *grpc.Server, i *Interceptor) (err error) {
	infos := s.GetServiceInfo()
	for sn, info := range infos {
		for _, m := range info.Methods {
			t := handlerType(m.IsClientStream, m.IsServerStream)

			for c := uint32(0); c <= 15; c++ {
				requestLabels := prometheus.Labels{
					"service": sn,
					"handler": m.Name,
					"code":    codes.Code(c).String(),
					"type":    t,
				}
				messageLabels := prometheus.Labels{
					"service": sn,
					"handler": m.Name,
				}

				// client
				if _, err = i.monitoring.client.errors.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.client.requests.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.client.requestDuration.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.client.messagesReceived.GetMetricWith(messageLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.client.messagesSend.GetMetricWith(messageLabels); err != nil {
					return err
				}
				// server
				if _, err = i.monitoring.server.errors.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.server.requests.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.server.requestDuration.GetMetricWith(requestLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.server.messagesReceived.GetMetricWith(messageLabels); err != nil {
					return err
				}
				if _, err = i.monitoring.server.messagesSend.GetMetricWith(messageLabels); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Interceptor ...
type Interceptor struct {
	monitoring *monitoring
}

// NewInterceptor ...
func NewInterceptor() *Interceptor {
	return &Interceptor{
		monitoring: initMonitoring(),
	}
}

// Dialer ...
func (i *Interceptor) Dialer(f func(string, time.Duration) (net.Conn, error)) func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		i.monitoring.dialer.WithLabelValues(addr).Inc()
		return f(addr, timeout)
	}
}

// UnaryClient ...
func (i *Interceptor) UnaryClient() grpc.UnaryClientInterceptor {
	monitor := i.monitoring.client

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)
		code := grpc.Code(err)
		service, method := split(method)
		labels := prometheus.Labels{
			"service": service,
			"handler": method,
			"code":    code.String(),
			"type":    "unary",
		}
		if err != nil && code != codes.OK {
			monitor.errors.With(labels).Add(1)
		}

		elapsed := float64(time.Since(start)) / float64(time.Microsecond)
		monitor.requestDuration.With(labels).Observe(elapsed)
		monitor.requests.With(labels).Add(1)

		return err
	}
}

// StreamClient ...
func (i *Interceptor) StreamClient() grpc.StreamClientInterceptor {
	monitor := i.monitoring.client

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()

		client, err := streamer(ctx, desc, cc, method, opts...)
		code := grpc.Code(err)
		service, method := split(method)
		labels := prometheus.Labels{
			"service": service,
			"handler": method,
			"code":    code.String(),
			"type":    handlerType(desc.ClientStreams, desc.ServerStreams),
		}
		if err != nil && code != codes.OK {
			monitor.errors.With(labels).Add(1)
		}

		elapsed := float64(time.Since(start)) / float64(time.Microsecond)
		monitor.requestDuration.With(labels).Observe(elapsed)
		monitor.requests.With(labels).Add(1)

		return &monitoredClientStream{ClientStream: client, monitor: monitor, labels: prometheus.Labels{
			"service": service,
			"handler": method,
		}}, nil
	}
}

// UnaryServer ...
func (i *Interceptor) UnaryServer() grpc.UnaryServerInterceptor {
	monitor := i.monitoring.server

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, err := handler(ctx, req)
		code := grpc.Code(err)
		service, method := split(info.FullMethod)
		labels := prometheus.Labels{
			"service": service,
			"handler": method,
			"code":    code.String(),
			"type":    "unary",
		}
		if err != nil && code != codes.OK {
			monitor.errors.With(labels).Add(1)
		}

		elapsed := float64(time.Since(start)) / float64(time.Microsecond)
		monitor.requestDuration.With(labels).Observe(elapsed)
		monitor.requests.With(labels).Add(1)

		return res, err
	}
}

// StreamServer ...
func (i *Interceptor) StreamServer() grpc.StreamServerInterceptor {
	monitor := i.monitoring.server

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		service, method := split(info.FullMethod)
		err := handler(srv, &monitoredServerStream{ServerStream: ss, labels: prometheus.Labels{
			"service": service,
			"handler": method,
		}, monitor: monitor})
		code := grpc.Code(err)
		labels := prometheus.Labels{
			"service": service,
			"handler": method,
			"code":    code.String(),
			"type":    handlerType(info.IsClientStream, info.IsServerStream),
		}
		if err != nil && code != codes.OK {
			monitor.errors.With(labels).Add(1)
		}

		elapsed := float64(time.Since(start)) / float64(time.Microsecond)
		monitor.requestDuration.With(labels).Observe(elapsed)
		monitor.requests.With(labels).Add(1)

		return err
	}
}

// ServerRegistry returns the prometheus registry with all server metrics
func (i *Interceptor) ServerRegistry() *prometheus.Registry {
	return i.monitoring.serverRegistry
}

// ClientRegistry returns the prometheus registry with all client metrics
func (i *Interceptor) ClientRegistry() *prometheus.Registry {
	return i.monitoring.clientRegistry
}

type monitoring struct {
	dialer         *prometheus.CounterVec
	server         *monitor
	client         *monitor
	serverRegistry *prometheus.Registry
	clientRegistry *prometheus.Registry
}

type monitor struct {
	requests         *prometheus.CounterVec
	requestDuration  *prometheus.SummaryVec
	messagesReceived *prometheus.CounterVec
	messagesSend     *prometheus.CounterVec
	errors           *prometheus.CounterVec
}

func initMonitoring() *monitoring {
	dialer := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "reconnects_total",
			Help:      "Total number of reconnects made by client.",
		},
		[]string{"address"},
	)
	serverRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Total number of RPC requests received by server.",
		},
		[]string{"service", "handler", "code", "type"},
	)
	serverReceivedMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "received_messages_total",
			Help:      "Total number of RPC messages received by server.",
		},
		[]string{"service", "handler"},
	)
	serverSendMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "send_messages_total",
			Help:      "Total number of RPC messages send by server.",
		},
		[]string{"service", "handler"},
	)
	serverRequestDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "request_duration_microseconds",
			Help:      "The RPC request latencies in microseconds on server side.",
		},
		[]string{"service", "handler", "code", "type"},
	)
	serverErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "errors_total",
			Help:      "Total number of errors that happen during RPC calles on server side.",
		},
		[]string{"service", "handler", "code", "type"},
	)

	serverRegistry := prometheus.NewRegistry()

	serverRegistry.MustRegister(dialer)
	serverRegistry.MustRegister(serverRequests)
	serverRegistry.MustRegister(serverRequestDuration)
	serverRegistry.MustRegister(serverReceivedMessages)
	serverRegistry.MustRegister(serverSendMessages)
	serverRegistry.MustRegister(serverErrors)

	clientRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "requests_total",
			Help:      "Total number of RPC requests made by client.",
		},
		[]string{"service", "handler", "code", "type"},
	)
	clientReceivedMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "received_messages_total",
			Help:      "Total number of RPC messages received.",
		},
		[]string{"service", "handler"},
	)
	clientSendMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "send_messages_total",
			Help:      "Total number of RPC messages send.",
		},
		[]string{"service", "handler"},
	)
	clientRequestDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "request_duration_microseconds",
			Help:      "The RPC request latencies in microseconds.",
		},
		[]string{"service", "handler", "code", "type"},
	)
	clientErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Subsystem: "client",
			Name:      "errors_total",
			Help:      "Total number of errors that happen during RPC calles.",
		},
		[]string{"service", "handler", "code", "type"},
	)

	clientRegistry := prometheus.NewRegistry()

	clientRegistry.MustRegister(clientRequests)
	clientRegistry.MustRegister(clientRequestDuration)
	clientRegistry.MustRegister(clientReceivedMessages)
	clientRegistry.MustRegister(clientSendMessages)
	clientRegistry.MustRegister(clientErrors)

	return &monitoring{
		dialer:         dialer,
		serverRegistry: serverRegistry,
		server: &monitor{
			requests:         serverRequests,
			requestDuration:  serverRequestDuration,
			messagesReceived: serverReceivedMessages,
			messagesSend:     serverSendMessages,
			errors:           serverErrors,
		},
		clientRegistry: clientRegistry,
		client: &monitor{
			requests:         clientRequests,
			requestDuration:  clientRequestDuration,
			messagesReceived: clientReceivedMessages,
			messagesSend:     clientSendMessages,
			errors:           clientErrors,
		},
	}
}

type monitoredServerStream struct {
	grpc.ServerStream
	labels  prometheus.Labels
	monitor *monitor
}

func (mss *monitoredServerStream) SendMsg(m interface{}) error {
	err := mss.ServerStream.SendMsg(m)
	if err == nil {
		mss.monitor.messagesSend.With(mss.labels).Inc()
	}
	return err
}

func (mss *monitoredServerStream) RecvMsg(m interface{}) error {
	err := mss.ServerStream.RecvMsg(m)
	if err == nil {
		mss.monitor.messagesReceived.With(mss.labels).Inc()
	}
	return err
}

type monitoredClientStream struct {
	grpc.ClientStream
	labels  prometheus.Labels
	monitor *monitor
}

func (mcs *monitoredClientStream) SendMsg(m interface{}) error {
	err := mcs.ClientStream.SendMsg(m)
	if err == nil {
		mcs.monitor.messagesSend.With(mcs.labels).Inc()
	}
	return err
}

func (mcs *monitoredClientStream) RecvMsg(m interface{}) error {
	err := mcs.ClientStream.RecvMsg(m)
	if err == nil {
		mcs.monitor.messagesReceived.With(mcs.labels).Inc()
	}
	return err
}

func handlerType(clientStream, serverStream bool) string {
	switch {
	case !clientStream && !serverStream:
		return "unary"
	case !clientStream && serverStream:
		return "server_stream"
	case clientStream && !serverStream:
		return "client_stream"
	default:
		return "bidirectional_stream"
	}
}

func split(name string) (string, string) {
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[1:i], name[i+1:]
	}
	return "unknown", "unknown"
}
