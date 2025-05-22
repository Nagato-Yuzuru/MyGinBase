package logger

import (
	"GinBase/pkg/errs"
	"context"
	"errors"
	"os"
	"sync"
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

type ZapLogger struct {
	logger *zap.Logger
}

var (
	zapLogger     *ZapLogger
	zapLoggerOnce sync.Once
)

type loggerParas struct {
	level      string
	isJSON     bool
	appName    string
	appVersion string
}

func ProvideZapLogger(lp *loggerParas) (*ZapLogger, error) {
	if lp == nil {
		return nil, errs.WrapCodeError(errs.ErrNilPointer, errors.New("loggerParas is nil"))
	}

	level := lp.level
	isJSON := lp.isJSON
	appName := lp.appName
	appVersion := lp.appVersion

	zapLoggerOnce.Do(
		func() {
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

			bufferSize := 256 * 1024
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
			logger := zap.New(
				asyncCore,
				zap.AddCaller(),
				zap.AddCallerSkip(1),
				zap.AddStacktrace(zapcore.ErrorLevel),
			).With(fields...)

			zapLogger = &ZapLogger{logger: logger}
		},
	)

	if zapLogger == nil {
		return nil, errors.New("zap logger is not initialized, call zap logger init first")
	}

	return zapLogger, nil

}

func (l *ZapLogger) getTraceFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return []zap.Field{
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		}
	}
	return nil
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Debug(msg, allFields...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Info(msg, allFields...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Warn(msg, allFields...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Error(msg, allFields...)
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Fatal(msg, allFields...)
}

func (l *ZapLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getTraceFields(ctx), fields...)
	l.logger.Panic(msg, allFields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
}

// Sync flushes any buffered log entries.
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
