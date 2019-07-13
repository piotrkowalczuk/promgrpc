package promgrpc

type rpcTag struct {
	isFailFast      string
	service         string
	method          string
	clientUserAgent string
}
