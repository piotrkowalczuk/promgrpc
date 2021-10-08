package useragent_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/alexeyxo/promgrpc/v4/internal/useragent"
	"google.golang.org/grpc/stats"
)

func TestStore_ClientSide(t *testing.T) {
	const ua = "n/a/y"

	ctx := context.Background()
	req := &stats.Begin{Client: true}

	var store useragent.Store
	if res := store.ClientSide(ctx, req); res != ua {
		t.Fatalf("wrong result: %s", res)
	}

	for i := 0; i < 10; i++ {
		if res := store.ClientSide(ctx, req); res != ua {
			t.Fatalf("wrong result: %s", res)
		}
	}
}

func BenchmarkStore_ClientSide_notAvailableYet(b *testing.B) {
	const ua = "n/a/y"

	ctx := context.Background()
	req := &stats.Begin{Client: true}

	var store useragent.Store
	if res := store.ClientSide(ctx, req); res != ua {
		b.Fatalf("wrong result: %s", res)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if res := store.ClientSide(ctx, req); res != ua {
			b.Fatalf("wrong result: %s", res)
		}
	}
}

func BenchmarkStore_ClientSide_available(b *testing.B) {
	const ua = "user-agent-store-test"
	ctx := context.Background()
	req := &stats.OutHeader{Client: true, Header: metadata.MD{"user-agent": []string{ua}}}

	var store useragent.Store
	if res := store.ClientSide(ctx, req); res != ua {
		b.Fatalf("wrong result: %s", res)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if res := store.ClientSide(ctx, req); res != ua {
			b.Fatalf("wrong result: %s", res)
		}
	}
}
