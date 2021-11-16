package types

import (
	"go.uber.org/zap/zapcore"
)

func (t *TradeValue) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("price", t.Price.String())
	enc.AddString("volume", t.Size.String())
	return nil
}

func (t *TradeMsg) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("type", t.Type)
	enc.AddString("type", t.Message)
	enc.AddString("type", t.Reason)
	enc.AddFloat64("price", t.Price)
	enc.AddFloat64("volume", t.Size)
	return nil
}

func (v *VWAPResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("productID", v.ProductID)
	enc.AddString("vwap", v.Vwap.String())
	return nil
}
