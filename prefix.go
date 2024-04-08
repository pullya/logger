package logger

import (
	"context"
	"fmt"
)

func SetPrefix(ctx context.Context, prefix string) context.Context {
	if prefix == "" {
		return ctx
	}

	if parentPrefix := getPrefix(ctx); parentPrefix != "" {
		prefix = parentPrefix + prefix
	}

	return context.WithValue(ctx, ctxPrefixKeyVal, prefix+": ")
}

func getPrefix(ctx context.Context) string {
	val := ctx.Value(ctxPrefixKeyVal)

	if val == nil {
		return ""
	}

	prefix, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("unexpected context log prefix value type: %T", val))
	}

	return prefix
}

type ctxPrefixKeyType struct{}

var ctxPrefixKeyVal ctxPrefixKeyType = struct{}{}
