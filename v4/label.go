package promgrpc

import (
	"context"
	"reflect"

	"github.com/alexeyxo/promgrpc/v4/internal/useragent"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

type rpcTagLabels struct {
	isFailFast      string
	service         string
	method          string
	clientUserAgent string
	code            string
}

type connTagLabels struct {
	remoteAddr      string
	localAddr       string
	clientUserAgent string
}

// HandleRPCLabelFunc type represents a function signature that can be passed into a stats handler and used instead of default one.
// That way caller gets the ability to modify the way labels are assembled.
type HandleRPCLabelFunc func(context.Context, stats.RPCStats) []string

// TagRPCLabelFunc type represents a function signature that can be passed into StatsHandlerWithTagRPCLabelsFunc.
type TagRPCLabelFunc func(context.Context, *stats.RPCTagInfo) context.Context

const structTag = "promgrpc"

type supportedLabels struct {
	// keep alphabetical order
	ClientUserAgent bool `promgrpc:"grpc_client_user_agent"`
	Code            bool `promgrpc:"grpc_code"`
	IsFailFast      bool `promgrpc:"grpc_is_fail_fast"`
	LocalAddr       bool `promgrpc:"grpc_local_addr"`
	Method          bool `promgrpc:"grpc_method"`
	RemoteAddr      bool `promgrpc:"grpc_remote_addr"`
	Service         bool `promgrpc:"grpc_service"`
}

func (l supportedLabels) labels() []string {
	var res []string

	v := reflect.ValueOf(l)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Bool() {
			res = append(res, v.Type().Field(i).Tag.Get(structTag))
		}
	}

	return res
}

func (l supportedLabels) isKnown(name string) bool {
	v := reflect.ValueOf(l)

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get(structTag) == name {
			return true
		}
	}

	return false
}

func (l *supportedLabels) enable(name string) {
	v := reflect.ValueOf(l).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Tag.Get(structTag) == name {
			v.Field(i).SetBool(true)
		}
	}
}

func (l *supportedLabels) subsetConn(tag connTagLabels) []string {
	res := make([]string, 0, 3)

	if l.ClientUserAgent {
		res = append(res, tag.clientUserAgent)
	}
	if l.LocalAddr {
		res = append(res, tag.localAddr)
	}
	if l.RemoteAddr {
		res = append(res, tag.remoteAddr)
	}

	return res
}

func (l *supportedLabels) subsetRPC(tag rpcTagLabels) []string {
	res := make([]string, 0, 5)

	if l.ClientUserAgent {
		res = append(res, tag.clientUserAgent)
	}
	if l.Code {
		res = append(res, tag.code)
	}
	if l.IsFailFast {
		res = append(res, tag.isFailFast)
	}
	if l.Method {
		res = append(res, tag.method)
	}
	if l.Service {
		res = append(res, tag.service)
	}

	return res
}

type clientSideLabelsHandler struct {
	supportedLabels supportedLabels
	userAgentStore  useragent.Store
}

func (h *clientSideLabelsHandler) labelsTagConn(ctx context.Context) []string {
	tag := ctx.Value(tagConnKey).(connTagLabels)

	return h.supportedLabels.subsetConn(tag)
}

func (h *clientSideLabelsHandler) labelsTagRPC(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)

	if h.supportedLabels.ClientUserAgent {
		tag.clientUserAgent = h.userAgentStore.ClientSide(ctx, stat)
	}
	if h.supportedLabels.Code {
		tag.code = status.Code(stat.(*stats.End).Error).String()
	}

	return h.supportedLabels.subsetRPC(tag)
}

type serverSideLabelsHandler struct {
	supportedLabels supportedLabels
}

func (h *serverSideLabelsHandler) labelsTagConn(ctx context.Context) []string {
	tag := ctx.Value(tagConnKey).(connTagLabels)

	return h.supportedLabels.subsetConn(tag)
}

func (h *serverSideLabelsHandler) labelsTagRPC(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTagLabels)
	if h.supportedLabels.Code {
		tag.code = status.Code(stat.(*stats.End).Error).String()
	}

	return h.supportedLabels.subsetRPC(tag)
}
