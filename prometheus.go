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

// Interceptor ...
type Interceptor struct {
	monitoring *monitoring
}

// NewInterceptor ...
func NewInterceptor(labels prometheus.Labels) *Interceptor {
	return &Interceptor{
		monitoring: initMonitoring(labels),
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
			"service":     service,
			"method":      method,
			"code":        code.String(),
			"method_type": "unary",
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
			"service":     service,
			"method":      method,
			"code":        code.String(),
			"method_type": handlerType(desc.ClientStreams, desc.ServerStreams),
		}
		if err != nil && code != codes.OK {
			monitor.errors.With(labels).Add(1)
		}

		elapsed := float64(time.Since(start)) / float64(time.Microsecond)
		monitor.requestDuration.With(labels).Observe(elapsed)
		monitor.requests.With(labels).Add(1)

		return &monitoredClientStream{ClientStream: client, monitor: monitor, labels: prometheus.Labels{
			"service": service,
			"method":  method,
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
			"service":      service,
			"handler":      method,
			"code":         code.String(),
			"handler_type": "unary",
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
			"service":      service,
			"handler":      method,
			"code":         code.String(),
			"handler_type": handlerType(info.IsClientStream, info.IsServerStream),
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

type monitoring struct {
	dialer *prometheus.CounterVec
	server *monitor
	client *monitor
}

type monitor struct {
	requests         *prometheus.CounterVec
	requestDuration  *prometheus.SummaryVec
	messagesReceived *prometheus.CounterVec
	messagesSend     *prometheus.CounterVec
	errors           *prometheus.CounterVec
}

func initMonitoring(constLabels prometheus.Labels) *monitoring {
	dialer := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "reconnects_total",
			Help:        "Total number of reconnects made by client.",
			ConstLabels: constLabels,
		},
		[]string{"address"},
	)
	serverRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "server",
			Name:        "requests_total",
			Help:        "Total number of RPC requests received by server.",
			ConstLabels: constLabels,
		},
		[]string{"service", "handler", "code", "handler_type"},
	)
	serverReceivedMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "server",
			Name:        "received_messages_total",
			Help:        "Total number of RPC messages received by server.",
			ConstLabels: constLabels,
		},
		[]string{"service", "handler"},
	)
	serverSendMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "server",
			Name:        "send_messages_total",
			Help:        "Total number of RPC messages send by server.",
			ConstLabels: constLabels,
		},
		[]string{"service", "handler"},
	)
	serverRequestDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:   "grpc",
			Subsystem:   "server",
			Name:        "request_duration_microseconds",
			Help:        "The RPC request latencies in microseconds on server side.",
			ConstLabels: constLabels,
		},
		[]string{"service", "handler", "code", "handler_type"},
	)
	serverErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "server",
			Name:        "errors_total",
			Help:        "Total number of errors that happen during RPC calles on server side.",
			ConstLabels: constLabels,
		},
		[]string{"service", "handler", "code", "handler_type"},
	)

	// TODO: re-implement for prometheus v0.9.0
	dialer = prometheus.MustRegisterOrGet(dialer).(*prometheus.CounterVec)

	// TODO: re-implement for prometheus v0.9.0
	serverRequests = prometheus.MustRegisterOrGet(serverRequests).(*prometheus.CounterVec)
	serverRequestDuration = prometheus.MustRegisterOrGet(serverRequestDuration).(*prometheus.SummaryVec)
	serverReceivedMessages = prometheus.MustRegisterOrGet(serverReceivedMessages).(*prometheus.CounterVec)
	serverSendMessages = prometheus.MustRegisterOrGet(serverSendMessages).(*prometheus.CounterVec)
	serverErrors = prometheus.MustRegisterOrGet(serverErrors).(*prometheus.CounterVec)

	clientRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "requests_total",
			Help:        "Total number of RPC requests made by client.",
			ConstLabels: constLabels,
		},
		[]string{"service", "method", "code", "method_type"},
	)
	clientReceivedMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "received_messages_total",
			Help:        "Total number of RPC messages received.",
			ConstLabels: constLabels,
		},
		[]string{"service", "method"},
	)
	clientSendMessages := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "send_messages_total",
			Help:        "Total number of RPC messages send.",
			ConstLabels: constLabels,
		},
		[]string{"service", "method"},
	)
	clientRequestDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "request_duration_microseconds",
			Help:        "The RPC request latencies in microseconds.",
			ConstLabels: constLabels,
		},
		[]string{"service", "method", "code", "method_type"},
	)
	clientErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "grpc",
			Subsystem:   "client",
			Name:        "errors_total",
			Help:        "Total number of errors that happen during RPC calles.",
			ConstLabels: constLabels,
		},
		[]string{"service", "method", "code", "method_type"},
	)

	// TODO: re-implement for prometheus v0.9.0
	clientRequests = prometheus.MustRegisterOrGet(clientRequests).(*prometheus.CounterVec)
	clientRequestDuration = prometheus.MustRegisterOrGet(clientRequestDuration).(*prometheus.SummaryVec)
	clientReceivedMessages = prometheus.MustRegisterOrGet(clientReceivedMessages).(*prometheus.CounterVec)
	clientSendMessages = prometheus.MustRegisterOrGet(clientSendMessages).(*prometheus.CounterVec)
	clientErrors = prometheus.MustRegisterOrGet(clientErrors).(*prometheus.CounterVec)

	return &monitoring{
		dialer: dialer,
		server: &monitor{
			requests:         serverRequests,
			requestDuration:  serverRequestDuration,
			messagesReceived: serverReceivedMessages,
			messagesSend:     serverSendMessages,
			errors:           serverErrors,
		},
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
