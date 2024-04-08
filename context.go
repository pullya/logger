package logger

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	contextKeySender     = contextKey("sender")
	contextKeyEntrypoint = contextKey("entrypoint")
)

func Sender(ctx context.Context) (string, bool) {
	valueStr, ok := ctx.Value(contextKeySender).(string)
	return valueStr, ok
}

func SetSender(ctx context.Context, sender string) context.Context {
	return context.WithValue(ctx, contextKeySender, sender)
}

func Entrypoint(ctx context.Context) (string, bool) {
	valueStr, ok := ctx.Value(contextKeyEntrypoint).(string)
	return valueStr, ok
}

func SetEntrypoint(ctx context.Context, entrypoint string) context.Context {
	return context.WithValue(ctx, contextKeyEntrypoint, entrypoint)
}
