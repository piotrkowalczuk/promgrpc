package promgrpc

import "context"

type ctxKeyType byte

const ctxKeyDynamicLabelValues ctxKeyType = iota

func DynamicLabelValuesFromCtx(ctx context.Context) map[string]string {
	return fromCtx[map[string]string](ctx, ctxKeyDynamicLabelValues)
}

func DynamicLabelValuesToCtx(ctx context.Context, dynamicLabelValues map[string]string) context.Context {
	return toCtx(ctx, ctxKeyDynamicLabelValues, dynamicLabelValues)
}

func toCtx[T any](ctx context.Context, key ctxKeyType, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func fromCtx[T any](ctx context.Context, key ctxKeyType) T {
	val := ctx.Value(key)
	if val == nil {
		tp := new(T)
		return *tp
	}
	return val.(T)
}
