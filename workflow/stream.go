// nolint:errcheck
package workflow

import (
	"context"
	"fmt"
	"math/big"
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
		productsVwap: vwap.New(cfg.ProductIDs, cfg.WindowsSize),
	}
}

// GetTradesQConsumer returns a trades queue consumer that receives
// trade values
func (c Client) GetTradesQConsumer() types.TradesQConsumer {
	return c.q
}

// TradesToVwap pipes trades to the go routines pool and receives back here the
// transformed VWAP results by the Results Queue.
func (c *Client) TradesToVwap(ctx context.Context) error {
	logger := log.FromContext(ctx)

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
	go ingestTradesStream(ctx, c.conn, c.GetTradesQConsumer(), doneTradesStreaming)

	return c.IngestVWAPResults(ctx, logger, doneTradesStreaming)
}

func (c *Client) IngestVWAPResults(ctx context.Context, logger *zap.Logger, doneTradesStreaming chan struct{}) error {
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

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping trades stream ingestion")
			return

		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Error("msg reading erred", zap.Error(err))
				return
			}

			msgType, idx := types.ParseType(msg)
			if idx == -1 {
				continue
			}

			switch msgType {
			case server.SubAckMsgType:
				logger.Info("Subscribed:")
			/*
			* Undocumented message type? Appears to propagate the same info as
			* `match`
			 */
			case server.MatchLastMsgType:
				fallthrough
			case server.MatchMsgType:
				msgProductID, idx := types.ParseProductID(msg)
				if idx == -1 {
					logger.Error("Failed to parse the product:"+string(msg))
					continue
				}

				msgPrice, idx := types.ParsePrice(msg)
				if idx == -1 {
					logger.Error("Failed to parse the price:"+string(msg))
					continue
				}

				msgVolume, idx := types.ParseVolume(msg)
				if idx == -1 {
					logger.Error("Failed to parse the volume:"+string(msg))
					continue
				}

				tradeValue := getMemPoolTradeVal()
				tradeValue.ProductID = msgProductID
				tradeValue.Price.SetFloat64(msgPrice)
				tradeValue.Size.SetFloat64(msgVolume)

				broadcast <- tradeValue

				logger.Debug(
					"received trade",
					zap.Object(tradeValue.ProductID, tradeValue),
				)
			case server.ErrorMsgType:
				logger.Error(
					"socket error",
					zap.String("msg", string(msg)),
				)
			default:
				logger.Warn(
					"unknown socket message",
					zap.String(msgType, string(msg)),
				)
			}
		}
	}
}

func getMemPoolTradeVal() *types.TradeValue {
	tradeValue := types.TradeValueMemPool.Get().(*types.TradeValue)
	tradeValue.Price = types.BigFloatMemPool.Get().(*big.Float)
	tradeValue.Size = types.BigFloatMemPool.Get().(*big.Float)
	return tradeValue
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
