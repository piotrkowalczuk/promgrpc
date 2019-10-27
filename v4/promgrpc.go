package promgrpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	namespace            = "grpc"
	labelService         = "grpc_service"
	labelMethod          = "grpc_method"
	labelCode            = "grpc_code"
	labelIsFailFast      = "grpc_is_fail_fast"
	labelRemoteAddr      = "grpc_remote_addr"
	labelLocalAddr       = "grpc_local_addr"
	labelClientUserAgent = "grpc_client_user_agent"
)

type ctxKey int

var (
	tagRPCKey  ctxKey = 1
	tagConnKey ctxKey = 3
)

func split(name string) (string, string) {
	if i := strings.LastIndex(name, "/"); i >= 0 {
		return name[1:i], name[i+1:]
	}
	return "unknown", "unknown"
}

func userAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua, ok := md["user-agent"]; ok {
			return ua[0]
		}
	}
	return "n/a"
}
