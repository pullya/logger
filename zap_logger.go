package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ServiceInfo struct {
	Environment string
	Name        string
	Version     string
	InstanceID  string
}

// ZapLogger is a wrapper for Zap to add context variables
type ZapLogger struct {
	zapLogger *zap.SugaredLogger
}

var (
	defaultLogger *ZapLogger
	defaultLevel  = zap.NewAtomicLevelAt(zap.ErrorLevel)
)

func init() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	loggerConfig.Level = defaultLevel

	logger, err := loggerConfig.Build()
	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		return
	}

	zap.ReplaceGlobals(logger)
	defaultLogger = &ZapLogger{zapLogger: withDefaultFields(logger.Sugar(), nil)}
}

func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

func Level() zapcore.Level {
	return defaultLevel.Level()
}

func SetDefaultFields(serviceInfo *ServiceInfo) {
	defaultLogger.zapLogger = withDefaultFields(defaultLogger.zapLogger, serviceInfo)
}

func NewZapLogger(logLevel LogLevel, serviceInfo *ServiceInfo) (*ZapLogger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	zapLogLevel, err := zap.ParseAtomicLevel(string(logLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}
	loggerConfig.Level = zapLogLevel

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return &ZapLogger{
		zapLogger: withDefaultFields(logger.Sugar(), serviceInfo),
	}, nil
}

func (l *ZapLogger) withContext(ctx context.Context) *zap.SugaredLogger {
	var (
		fields        = getFields(ctx)
		sender, _     = Sender(ctx)
		entrypoint, _ = Entrypoint(ctx)
	)

	if fields == nil {
		fields = make(Fields)
	}

	for k, v := range map[string]string{
		"sender":     sender,
		"entrypoint": entrypoint,

		// log trace and span ids according to https://b2bdevpro.atlassian.net/wiki/spaces/UPLATFORM/pages/380174341
		"trace.id": getTraceID(ctx),
		"span.id":  getSpanID(ctx),
	} {
		if v != "" {
			fields[k] = v
		}
	}

	return l.zapLogger.With(sugarFields(fields)...)
}

func prepMsg(ctx context.Context, msg string) string {
	if prefix := getPrefix(ctx); prefix != "" {
		return prefix + msg
	}
	return msg
}

// prepArgs converts log args to slice with message and zap.String fields
func prepArgs(ctx context.Context, args ...any) []any {
	var variables []any
	if len(args) != 0 {
		variables = append(variables, fmt.Sprintf(prepMsg(ctx, "%v"), args[0]))
	}

	if len(args) > 1 {
		i := 1
		for i < len(args) {
			var k, v string
			if i%2 != 0 {
				k = fmt.Sprintf("%v", args[i])
			}
			if i+1 < len(args) {
				v = fmt.Sprintf("%v", args[i+1])
			}

			variables = append(variables, zap.String(k, v))
			i++
		}
	}

	return variables
}

func (l *ZapLogger) Fatal(ctx context.Context, args ...any) {
	l.withContext(ctx).Fatal(prepArgs(ctx, args)...)
}

func (l *ZapLogger) Fatalf(ctx context.Context, template string, args ...any) {
	l.withContext(ctx).Fatalf(template, args...)
}

func (l *ZapLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	l.withContext(ctx).Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(ctx context.Context, args ...any) {
	l.withContext(ctx).Error(prepArgs(ctx, args)...)
}

func (l *ZapLogger) Errorf(ctx context.Context, template string, args ...any) {
	l.withContext(ctx).Errorf(template, args...)
}

func (l *ZapLogger) Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	l.withContext(ctx).Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(ctx context.Context, args ...any) {
	l.withContext(ctx).Info(prepArgs(ctx, args)...)
}

func (l *ZapLogger) Infof(ctx context.Context, template string, args ...any) {
	l.withContext(ctx).Infof(template, args...)
}

func (l *ZapLogger) Infow(ctx context.Context, msg string, keysAndValues ...any) {
	l.withContext(ctx).Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(ctx context.Context, args ...any) {
	l.withContext(ctx).Warn(prepArgs(ctx, args)...)
}

func (l *ZapLogger) Warnf(ctx context.Context, template string, args ...any) {
	l.withContext(ctx).Warnf(template, args...)
}

func (l *ZapLogger) Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	l.withContext(ctx).Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Debug(ctx context.Context, args ...any) {
	l.withContext(ctx).Debug(prepArgs(ctx, args)...)
}

func (l *ZapLogger) Debugf(ctx context.Context, template string, args ...any) {
	l.withContext(ctx).Debugf(template, args...)
}

func (l *ZapLogger) Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	l.withContext(ctx).Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) With(args ...any) Logger {
	return &ZapLogger{
		zapLogger: l.zapLogger.With(args...),
	}
}

func withDefaultFields(logger *zap.SugaredLogger, serviceInfo *ServiceInfo) *zap.SugaredLogger {
	if serviceInfo == nil {
		return logger
	}

	var kvs []any
	for k, v := range map[string]string{
		// log service information according to https://b2bdevpro.atlassian.net/wiki/spaces/UPLATFORM/pages/380174341
		"service.environment": serviceInfo.Environment,
		"service.name":        serviceInfo.Name,
		"service.version":     serviceInfo.Version,
		"service.instance.id": serviceInfo.InstanceID,
	} {
		if v != "" {
			kvs = append(kvs, k, v)
		}
	}

	if len(kvs) > 0 {
		return logger.With(kvs...)
	}

	return logger
}
