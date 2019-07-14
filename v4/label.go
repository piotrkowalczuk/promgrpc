package promgrpc

import (
	"context"

	"google.golang.org/grpc/stats"
)

type rpcTag struct {
	isFailFast      string
	service         string
	method          string
	clientUserAgent string
}

// RPCLabelFunc ...
type RPCLabelFunc func(ctx context.Context, stat stats.RPCStats) []string
