package promgrpc

import (
	"context"

	"google.golang.org/grpc/stats"
)

const (
	labelService         = "grpc_service"
	labelMethod          = "grpc_method"
	labelCode            = "grpc_code"
	labelIsFailFast      = "grpc_is_fail_fast"
	labelRemoteAddr      = "grpc_remote_addr"
	labelLocalAddr       = "grpc_local_addr"
	labelClientUserAgent = "grpc_client_user_agent"
)

type rpcTagLabels struct {
	isFailFast      string
	service         string
	method          string
	clientUserAgent string
}

type connTagLabels struct {
	remoteAddr      string
	localAddr       string
	clientUserAgent string
}

// HandleRPCLabelFunc type represents a function signature that can be passed into a stats handler and used instead of default one.
// That way caller gets the ability to modify the way labels are assembled.
type HandleRPCLabelFunc func(context.Context, stats.RPCStats) []string

type TagRPCLabelFunc func(context.Context, *stats.RPCTagInfo) context.Context

type AdditionalLabelValuesFunc func(ctx context.Context) []string
