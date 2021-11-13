// nolint:errcheck
package workflow

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/blewater/zh/cmd"
	"github.com/blewater/zh/log"
	"github.com/blewater/zh/server"
	"github.com/blewater/zh/types"
	"github.com/blewater/zh/vwap"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Client to socket listen -> ingestTradesStream trades
type Client struct {
	cfg cmd.Config

	productsVwap *vwap.ProductsVwap

	// Inbound messages to be processed
	q chan *types.TradeValue

	// Inbound connection
	conn *websocket.Conn
}

func New(cfg cmd.Config) Client {
	return Client{
		q:            make(types.TradesQ, cfg.WorkerPoolSize),
		cfg:          cfg,
		productsVwap: vwap.New(cfg.ProductIDs, cfg.WorkerPoolSize),
	}
}

// GetTradesQProducer returns a trades generator or producer that produces trade
// values
func (c Client) GetTradesQProducer() types.TradesQProducer {
	return c.q
}

// GetTradesQConsumer returns a trades queue consumer that receives
// trade values
func (c Client) GetTradesQConsumer() types.TradesQConsumer {
	return c.q
}

func (c *Client) PipeTradesToVwapQ(ctx context.Context) error {
	defer close(c.q)

	logger := log.FromContext(ctx)
	defer logger.Sync()

	var err error
	c.conn, err = server.Connect(ctx, c.cfg.SocketURL)
	if err != nil {
		return err
	}
	defer c.conn.Close()

	if err := server.Subscribe(ctx, c.conn, c.cfg.ProductIDs); err != nil {
		return err
	}

	doneTradesStreaming := make(chan struct{})
	go ingestTradesStream(ctx, c.conn, c.q, doneTradesStreaming)

	for {
		select {
		case res := <-c.productsVwap.GetResultsQ():
			_, _ = fmt.Fprintf(
				os.Stderr, "ProductID:%s VWAP:%f\n", res.ProductID, res.Vwap,
			)
			// recycle into the mem pool
			types.VWAPResultMemPool.Put(res)
		case <-ctx.Done():
			c.gracefulSocketClose(logger, doneTradesStreaming)
			return nil
		}
	}
}

func ingestTradesStream(ctx context.Context, conn *websocket.Conn, broadcast chan<- *types.TradeValue, quit chan<- struct{}) {
	defer close(quit)

	logger := log.FromContext(ctx)
	defer logger.Sync()

	var tradeMsg types.TradeMsg
	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping trades stream ingestion")
			return

		default:
			if err := conn.ReadJSON(&tradeMsg); err != nil {
				logger.Error("msg reading erred", zap.Error(err))
				return
			}
			switch tradeMsg.Type {
			case server.SubAckMsgType:
				logger.Info("Subscribed:")
			case server.MatchLastMsgType:
				fallthrough
			case server.MatchMsgType:
				tradeValue := &types.TradeValue{
					ProductID: tradeMsg.ProductID,
				}
				tradeValue.Price = big.NewFloat(0)
				tradeValue.Size = big.NewFloat(0)
				tradeValue.Price.Set(tradeMsg.Price)
				tradeValue.Size.Set(tradeMsg.Size)
				broadcast <- tradeValue

				logger.Debug("received trade", zap.Any("ticker", tradeValue))
			case server.ErrorMsgType:
				logger.Error(
					"socket error",
					zap.String(tradeMsg.Reason, tradeMsg.Message),
				)
			default:
				logger.Warn(
					"unknown socket message",
					zap.Any(tradeMsg.Type, tradeMsg),
				)
			}
		}
	}
}

func (c *Client) gracefulSocketClose(logger *zap.Logger, doneTradesStreaming <-chan struct{}) {
	defer logger.Sync()
	logger.Info("Closing socket")

	// Cleanly close the inboundConn by sending a close message and then
	// wait (with timeout) for the server to close the inboundConn.
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		logger.Error("write close error:", zap.Error(err))
	}

	select {
	case <-doneTradesStreaming:
	case <-time.After(time.Second):
	}
}