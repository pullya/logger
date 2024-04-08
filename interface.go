package logger

//go:generate go run go.uber.org/mock/mockgen@latest -destination mocks/Logger.go . Logger

import "context"

type Logger interface {
	Fatal(ctx context.Context, args ...any)
	Fatalf(ctx context.Context, template string, args ...any)
	Fatalw(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, args ...any)
	Errorf(ctx context.Context, template string, args ...any)
	Errorw(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, args ...any)
	Infof(ctx context.Context, template string, args ...any)
	Infow(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, keysAndValues ...any)
	Warnf(ctx context.Context, template string, args ...any)
	Warnw(ctx context.Context, msg string, keysAndValues ...any)
	Debug(ctx context.Context, args ...any)
	Debugf(ctx context.Context, template string, args ...any)
	Debugw(ctx context.Context, msg string, keysAndValues ...any)
	With(args ...any) Logger
}
