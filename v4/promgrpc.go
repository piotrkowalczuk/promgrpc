package promgrpc

//go:generate stringer -type=Subsystem -output=subsystem.stringer.go

import (
	"context"
	"strings"

	_ "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
)

const (
	namespace    = "grpc"
	labelService = "grpc_service"
	labelMethod  = "grpc_method"
	//labelType      = "grpc_type"
	//labelCode      = "grpc_code"
	//labelUserAgent = "grpc_user_agent"
	labelFailFast        = "grpc_fail_fast"
	labelRemoteAddr      = "grpc_remote_addr"
	labelLocalAddr       = "grpc_local_addr"
	labelClientUserAgent = "grpc_client_user_agent"
	labelServerUserAgent = "grpc_server_user_agent"
)

type ctxKey int

var (
	tagRPCKey  ctxKey = 1
	tagConnKey ctxKey = 2
)

type Subsystem int

var (
	Server Subsystem = 1
	Client Subsystem = 2
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
	return "not-set"
}
