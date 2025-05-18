package logger

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for the logger.
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new zap logger.
// It allows configuring the log level and whether to output in JSON format.
func NewZapLogger(level string, isJSON bool, appName string, appVersion string) (Logger, error) {
	logLevel := zapcore.InfoLevel
	if err := logLevel.Set(level); err != nil {
		logLevel = zapcore.InfoLevel // Default to Info if parsing fails
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if isJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 使用缓冲写入器来提高性能
	bufferSize := 256 * 1024 // 例如，256KB的缓冲区
	flushInterval := 5 * time.Second
	writer := zapcore.AddSync(os.Stdout)

	bufferedWriter := &zapcore.BufferedWriteSyncer{
		WS:            writer,
		Size:          bufferSize,
		FlushInterval: flushInterval,
	}

	asyncCore := zapcore.NewCore(
		encoder,
		bufferedWriter,
		logLevel,
	)

	// Add common fields
	fields := []zap.Field{
		zap.String("service.name", appName),
		zap.String("service.version", appVersion),
	}

	// Wrap the asyncCore with zap.New to create the logger
	logger := zap.New(asyncCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).With(fields...)

	return &zapLogger{logger: logger}, nil
}

func (l *zapLogger) getTraceFields(ctx context.Context) []zap.Field {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return []zap.Field{
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		}
	}
	return nil
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Debug(msg, allFields...)
}

func (l *zapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Info(msg, allFields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Warn(msg, allFields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Error(msg, allFields...)
}

func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Fatal(msg, allFields...)
}

func (l *zapLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Panic(msg, allFields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// Sync flushes any buffered log entries.
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
