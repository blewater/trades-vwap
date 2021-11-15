package workflow

import (
	"context"
	"testing"

	"github.com/blewater/zh/cmd"
	"github.com/blewater/zh/log"
	"github.com/blewater/zh/server"
	"go.uber.org/zap"
)

func TestClient_TradesToVwapQ(t *testing.T) {
	c, err := setupStream(5)
	if err != nil {
		return
	}
	defer c.conn.Close()

	t.Run(
		"trx", func(t *testing.T) {
			tradesToVwapTrxs(&c, 1)
		},
	)
}

func Benchmark_100_VWAP_Trx_1Thread(b *testing.B) {
	c, err := setupStream(1)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_2Threads(b *testing.B) {
	c, err := setupStream(2)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_3Threads(b *testing.B) {
	c, err := setupStream(3)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_5Threads(b *testing.B) {
	c, err := setupStream(5)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_10Threads(b *testing.B) {
	c, err := setupStream(10)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_100Threads(b *testing.B) {
	c, err := setupStream(100)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func Benchmark_100_VWAP_Trx_200Threads(b *testing.B) {
	c, err := setupStream(200)
	if err != nil {
		return
	}
	defer c.conn.Close()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tradesToVwapTrxs(&c, 100)
	}
}

func tradesToVwapTrxs(c *Client, trxCount int) {
	for i := 0; i < trxCount; {
		<-c.productsVwap.GetResultsQ()
		i++
	}
}

func setupStream(workersCnt uint16) (Client, error) {
	products := []string{"BTC-USD", "USDC-EUR", "ETH-BTC", "ETH-EUR", "BTC-EUR"}

	cfg := cmd.Config{
		WorkerPoolSize: workersCnt,
		WindowsSize:    200,
		DevLogLevel:    false,
		SocketURL:      "wss://ws-feed.exchange.coinbase.com",
		ProductIDs:     products,
	}

	w := New(cfg)

	logger := zap.NewNop()
	ctx := log.ContextWithLogger(context.Background(), logger)

	// nolint:errcheck
	go w.StartPool(ctx)

	var err error
	w.conn, err = server.Connect(ctx, w.cfg.SocketURL)
	if err != nil {
		return Client{}, err
	}

	if err := server.Subscribe(ctx, w.conn, products); err != nil {
		return Client{}, err
	}

	doneTradesStreaming := make(chan struct{})
	go ingestTradesStream(
		ctx, w.conn, w.GetTradesQConsumer(), doneTradesStreaming,
	)
	return w, nil
}
