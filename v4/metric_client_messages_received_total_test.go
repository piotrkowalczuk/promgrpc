package promgrpc_test

import (
	"context"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func TestNewClientMessagesReceivedTotalStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientMessagesReceivedTotalStatsHandler(promgrpc.NewClientMessagesReceivedTotalCounterVec()))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: true,
	})
	h.HandleRPC(ctx, &stats.InPayload{
		Client: false,
	})

	const metadata = `
		# HELP grpc_client_messages_received_total TODO
        # TYPE grpc_client_messages_received_total counter
	`
	expected := `
		grpc_client_messages_received_total{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_messages_received_total"); err != nil {
		t.Fatal(err)
	}
}

func TestClientMessagesReceivedTotalStatsHandler_HandleRPC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cases := map[string]func() (*promgrpc.StatsHandler, error){
		"simple": func() (*promgrpc.StatsHandler, error) {
			return promgrpc.NewStatsHandler(
				promgrpc.NewClientMessagesReceivedTotalStatsHandler(
					promgrpc.NewClientMessagesReceivedTotalCounterVec(),
				),
			), nil
		},
		"one-label": func() (*promgrpc.StatsHandler, error) {
			c := prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: "a",
					Subsystem: "b",
					Name:      "c",
					Help:      "d",
				},
				[]string{"grpc_is_fail_fast"},
			)

			return promgrpc.NewStatsHandler(
				promgrpc.NewClientMessagesReceivedTotalStatsHandler(
					c,
				),
			), nil
		},
		"one-label-three-curried": func() (*promgrpc.StatsHandler, error) {
			c := promgrpc.NewClientMessagesReceivedTotalCounterVec()

			c, err := c.CurryWith(prometheus.Labels{
				"grpc_client_user_agent": "curried",
				"grpc_service":           "curried",
				"grpc_method":            "curried",
			})
			if err != nil {
				return nil, err
			}
			return promgrpc.NewStatsHandler(
				promgrpc.NewClientMessagesReceivedTotalStatsHandler(
					c,
				),
			), nil
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			h, err := c()
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				if err := recover(); err != nil {
					t.Errorf("%s\n%s", err, string(debug.Stack()))
				}
			}()

			ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
				FullMethodName: "A/B",
				FailFast:       false,
			})
			h.HandleRPC(ctx, &stats.InPayload{
				Client: true,
				Data:   []byte("{}"),
			})
		})
	}
}

func TestNewClientMessagesReceivedTotalStatsHandler_panic(t *testing.T) {
	defer func() {
		if err := recover(); err != "metric partitioned with non-supported labels" {
			t.Errorf("wrong panic: %s", err)
		}
	}()

	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "a",
			Subsystem: "b",
			Name:      "c",
			Help:      "d",
		},
		[]string{"invalid_label"},
	)

	_ = promgrpc.NewClientMessagesReceivedTotalStatsHandler(c)
}
