package log

import (
	"context"

	"go.uber.org/zap"
)

type key int

var (
	loggerKey key
	logger    *zap.Logger
)

func New(devLogLevel bool, url string) *zap.Logger {
	var err error
	if devLogLevel {
		// human friendly output, lower panic threshold
		logger, err = zap.NewDevelopment(
			zap.AddCaller(),
			zap.Development(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
		// important to annotate the trades source
		logger = logger.Named(url)
	} else {
		// structured log for machine consumption
		logger, err = zap.NewProduction(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
		// important to annotate the trades source
		logger = logger.Named(url)
	}

	if err != nil {
		panic(err)
	}

	return logger
}

// ContextWithLogger returns a new Context that carries the logger value.
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns the logger stored in ctx, the package level as a fall-back
func FromContext(ctx context.Context) *zap.Logger {
	ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return logger
	}

	return ctxLogger
}
