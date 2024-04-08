package logger

import (
	"context"
	"fmt"
	"maps"
)

type Fields map[string]any

func SetFields(ctx context.Context, fields Fields) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	var (
		ctxFields = getFields(ctx)
		newFields = make(Fields, len(fields)+len(ctxFields))
	)

	maps.Copy(newFields, fields)
	maps.Copy(newFields, ctxFields)

	return context.WithValue(ctx, ctxFieldsKeyVal, newFields)
}

func getFields(ctx context.Context) Fields {
	val := ctx.Value(ctxFieldsKeyVal)

	if val == nil {
		return nil
	}

	f, ok := val.(Fields)
	if !ok {
		panic(fmt.Sprintf("unexpected context log fields value type: %T", val))
	}

	return f
}

func sugarFields(fields Fields) []any {
	if len(fields) == 0 {
		return nil
	}

	keysAndValues := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keysAndValues = append(keysAndValues, k, v)
	}

	return keysAndValues
}

type ctxFieldsKeyType struct{}

var ctxFieldsKeyVal ctxFieldsKeyType = struct{}{}
