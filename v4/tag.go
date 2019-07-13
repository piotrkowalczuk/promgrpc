package promgrpc

type rpcTag struct {
	isFailFast      bool
	service         string
	method          string
	clientUserAgent string
}
