package useragent

import (
	"context"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/stats"
)

const (
	notAvailable    = "n/a"
	notAvailableYet = "n/a/y"
)

// Store allows to access user-agent header concurrently and efficiently.
// Currently there is no way to retrieve user-agent during TagRPC stage:
// LINK: https://github.com/grpc/grpc-go/pull/3331
type Store struct {
	userAgent          string
	userAgentLock      sync.Mutex
	userAgentAvailable uint32
}

// ClientSide implements best effort logic to obtain user-agent.
// Works with grpc 1.28.0 and above.
func (s *Store) ClientSide(_ context.Context, stat stats.RPCStats) string {
	if !stat.IsClient() {
		return notAvailable
	}
	if atomic.LoadUint32(&s.userAgentAvailable) == 1 {
		if s.userAgent == "" {
			return notAvailable
		}
		return s.userAgent
	}

	if st, ok := stat.(*stats.OutHeader); ok {
		if ua, ok := st.Header["user-agent"]; ok && len(ua) == 1 {
			return s.store(ua[0])
		}
	}

	return notAvailableYet
}

func (s *Store) store(userAgent string) string {
	s.userAgentLock.Lock()
	defer s.userAgentLock.Unlock()
	if s.userAgentAvailable == 0 {
		s.userAgent = userAgent
		defer atomic.StoreUint32(&s.userAgentAvailable, 1)
	}

	if s.userAgent == "" {
		return notAvailableYet
	}
	return s.userAgent
}
