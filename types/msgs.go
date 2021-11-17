package types

import (
	"go.uber.org/zap/zapcore"
	"math/big"
	"sync"
)

// TradesQ is the queue of received trade values to be processed by the workers
// pool. Note about pointer use: Channels escape to heap so pointer objects are
// more efficient.
type TradesQ chan *TradeValue

// TradesQConsumer is the consumer alias for TradesQ
type TradesQConsumer chan<- *TradeValue

// SubReq {"type":"subscribe","product_ids":["BTC-USD","ETH-USD","ETH-BTC"],"channels":["matches"]}
type SubReq struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

var TradeValueMemPool = sync.Pool{
	New: func() interface{} {
		return new(TradeValue)
	},
}

// TradeValue represents the minimum data set to calculate the VWAP data points.
// e.g. "product_id":"ETH-USD","price":"4606.8","size":"0.00269988"
//  While more fields than defined are returned, these are
// the minimum required to meet the biz requirements of this service.
//
// e.g. a received match trade ticker
// {
// 	"type":"match","trade_id":178622422,
// 	"maker_order_id":"253c56b0-f115-4364-9e06-65ffd2412f3b",
// 	"taker_order_id":"928f8eb1-b6b4-4735-b12a-a512a0da684f","side":"sell",
// 	"size":"0.00269988","price":"4606.8","product_id":"ETH-USD",
// 	"sequence":22394045199,"time":"2021-11-10T21:37:07.988255Z"
// }
type TradeValue struct {
	ProductID string
	Price     *big.Float
	Size      *big.Float
}

type VWAPResult struct {
	ProductID string
	Vwap      *big.Float
}

type ResultsQ chan *VWAPResult

var BigFloatMemPool = sync.Pool{
	New: func() interface{} {
		return big.NewFloat(0)
	},
}

var VWAPResultMemPool = sync.Pool{
	New: func() interface{} {
		return new(VWAPResult)
	},
}

//
// Log marshalling methods to remove log reflection
//

func (t *TradeValue) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("price", t.Price.String())
	enc.AddString("volume", t.Size.String())
	return nil
}

func (v *VWAPResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("productID", v.ProductID)
	enc.AddString("vwap", v.Vwap.String())
	return nil
}

