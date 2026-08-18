package promgrpc

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"

	"google.golang.org/grpc/metadata"
)

const (
	namespace    = "grpc"
	notAvailable = "n/a"
)

type ctxKey int

const (
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

const (
	minMessageSizeBucket     = 32    // bytes
	minRequestDurationBucket = 0.001 // seconds
)

// exponentialBucketsRange returns cardinality exponential buckets covering
// [lower, limit]. The upper bucket is set to limit exactly, because
// prometheus.ExponentialBucketsRange accumulates rounding error and would
// otherwise report e.g. 29.999999999999947 for a 30s limit.
//
// A single bucket is returned for a degenerate range, that is when fewer than
// two buckets are requested or limit does not exceed lower. Both cases would
// otherwise produce equal or descending bounds, which panic once the histogram
// observes its first value.
func exponentialBucketsRange(lower, limit float64, cardinality int) []float64 {
	if cardinality < 2 || limit <= lower {
		return []float64{limit}
	}

	buckets := prometheus.ExponentialBucketsRange(lower, limit, cardinality)
	buckets[len(buckets)-1] = limit

	return buckets
}

func exponentialBucketsRangeForSize(limit float64, cardinality int) []float64 {
	if limit <= 0 {
		return prometheus.DefBuckets // for backward compatibility
	}

	return exponentialBucketsRange(minMessageSizeBucket, limit, cardinality)
}

func exponentialBucketsRangeForDuration(limit float64, cardinality int) []float64 {
	if limit <= 0 {
		return prometheus.DefBuckets // for backward compatibility
	}

	return exponentialBucketsRange(minRequestDurationBucket, limit, cardinality)
}
