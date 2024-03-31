package logging

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ShyunnY/jaeger-operator/internal/consts"
)

type Logger struct {
	level consts.LogLevel
	logr.Logger
}

// NewLogger 构建一个内部使用zap的logger
func NewLogger(level consts.LogLevel) Logger {

	logger := zapr.NewLogger(zap.New(initZapCore(level), zap.AddCaller()))
	return Logger{
		level:  level,
		Logger: logger,
	}
}

func DefaultLogger() Logger {
	return NewLogger(consts.LogLevelDebug)
}

func initZapCore(level consts.LogLevel) zapcore.Core {

	Level, err := zapcore.ParseLevel(string(level))
	if err != nil {
		Level = zap.InfoLevel
	}

	prodEncoderConfig := zap.NewProductionEncoderConfig()
	prodEncoderConfig.TimeKey = "ts"
	prodEncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(prodEncoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(Level),
	)

	return core
}

func (l Logger) WithValues(kvs ...any) Logger {
	l.Logger = l.Logger.WithValues(kvs...)
	return l
}

func (l Logger) WithName(name string) Logger {

	logger := NewLogger(l.level).Logger.WithName(name)
	return Logger{
		Logger: logger,
	}
}
