// nolint:errcheck
package main

import (
	"context"
	"github.com/blewater/zh/cmd"
	"github.com/blewater/zh/log"
	"github.com/blewater/zh/workflow"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

// nolint:errcheck
func main() {
	cfg, logger := bootstrap()

	w := workflow.New(cfg)

	ctx, cancel := context.WithCancel(
		log.ContextWithLogger(context.Background(), logger))

	// Precedes socket communication to start the pool
	go w.StartPool(ctx)

	go func() {
		if err := w.TradesToVwap(ctx); err != nil {
			// socket failed to connect
			os.Exit(1)
		}
	}()

	waitInterruptSignal(cancel)
}

func waitInterruptSignal(cancel context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	cancel()
}

func bootstrap() (cmd.Config, *zap.Logger) {
	cfg := cmd.Execute()

	logger := log.New(cfg.DevLogLevel, cfg.SocketURL)
	defer logger.Sync()

	logger.Info("Starting...")
	logger.Info("Subscribing",zap.Any("Pairs", cfg.ProductIDs))
	logger.Info("log-mode",zap.Any("development", cfg.DevLogLevel))
	logger.Info("Socket",zap.Any("URL", cfg.SocketURL))

	return cfg, logger
}