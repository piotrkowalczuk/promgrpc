package promgrpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/stats"

	"google.golang.org/grpc/metadata"
)

const (
	namespace    = "grpc"
	notAvailable = "n/a"
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

func userAgentOnServerSide(ctx context.Context, _ *stats.RPCTagInfo) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua, ok := md["user-agent"]; ok && len(ua) == 1 {
			return ua[0]
		}
	}
	return notAvailable
}
