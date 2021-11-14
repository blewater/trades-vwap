package workflow

import (
	"context"

	"github.com/blewater/zh/log"
	"github.com/blewater/zh/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// StartPool starts the configured number of go routines to crunch VWAP results
// streaming off the client queue.
func (c Client) StartPool(ctx context.Context) error {
	logger := log.FromContext(ctx)
	// nolint:errcheck
	defer logger.Sync()

	g, ctx := errgroup.WithContext(ctx)

	var w uint16
	for w = 1; w <= c.cfg.WorkerPoolSize; w++ {
		// localize to avoid capture
		w := w
		g.Go(
			func() error {
				for tradeValue := range c.q {
					if err := c.productsVwap.ProduceVwap(ctx,
						tradeValue.ProductID,
						tradeValue.Price,
						tradeValue.Size,
					); err != nil {
						logger.Error(tradeValue.ProductID, zap.Error(err))
					}

					logger.Debug("worker", zap.Uint16("ID", w))
					logger.Debug("received trade", zap.Any("trade", tradeValue))

					types.TradeValueMemPool.Put(tradeValue)

					// Note to reviewer. Returning err here would trigger cancel
					// and discontinue any further work in this pool because of
					// the err group. Likely for services the rationale
					// prohibits such action and the error is simply logged and
					// observed. Also cancelling the pool would require also to
					// call the service ctx.CancelFunc so that the stream
					// ingestion would cease too.
					// return err
				}
				return nil
			},
		)
	}

	return g.Wait()
}
