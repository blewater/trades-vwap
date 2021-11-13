package vwap

import (
	"context"
	"fmt"
	"github.com/blewater/zh/log"
	"go.uber.org/zap"
	"math/big"
	"sync"

	"github.com/blewater/zh/types"
)

// Cache of cumulative data points to calculate a fresh VWAP result for each
// new product price, volume pairs.
type vwapCache struct {
	TPV  *big.Float
	TVol *big.Float
	PV   *big.Float
	Vol  *big.Float
}

// Map of product->windowQueue of cached data points
type vwapMap map[string]*WindowQueue

// ProductsVwap is the container for calculating the queued results.
type ProductsVwap struct {
	sync.Mutex
	windowSize uint16
	vwapCache  vwapMap
	resultsQ   types.ResultsQ
}

var bigZero = big.NewFloat(0)

func New(productIDs []string, windowSize uint16) *ProductsVwap {
	prodVwap := &ProductsVwap{
		vwapCache:  make(vwapMap),
		resultsQ:   make(types.ResultsQ, windowSize),
		windowSize: windowSize,
	}
	for _, p := range productIDs {
		// Allocate capacity upfront
		prodVwap.vwapCache[p] = NewWindowQueue(windowSize)
	}

	return prodVwap
}

func (v *ProductsVwap) ProduceVwap(	ctx context.Context, productID string, price, volume *big.Float) error {
	logger := log.FromContext(ctx)
	// nolint:errcheck
	defer logger.Sync()

	newDataPoints := vwapCacheMemPool.Get().(*vwapCache)

	newDataPoints.TPV = big.NewFloat(0)
	newDataPoints.TVol = big.NewFloat(0)
	newDataPoints.PV = big.NewFloat(0)
	newDataPoints.Vol = big.NewFloat(0)

	newDataPoints.PV.Mul(price, volume)
	newDataPoints.Vol.Set(volume)

	v.Lock()

	window, ok := v.vwapCache[productID]
	// Cannot continue
	if !ok {
		v.Unlock()
		return fmt.Errorf("product ID %s not in VMAP newDataPoints", productID)
	}

	if window.len > 0 {
		prevDataPoints, ok := window.PeekLast()
		if !ok {
			v.Unlock()
			return fmt.Errorf(
				"could not access cached data set %d, %s", window.len, productID,
			)
		}

		// Add previous sums
		newDataPoints.TPV.Add(newDataPoints.PV, prevDataPoints.TPV)
		newDataPoints.TVol.Add(newDataPoints.Vol, prevDataPoints.TVol)
	} else {
		// first window item
		newDataPoints.TPV.Set(newDataPoints.PV)
		newDataPoints.TVol.Set(newDataPoints.Vol)
	}

	// drop window data point to make room for the new
	if window.len == v.windowSize {
		droppedDataPoints, ok := window.Pop()
		if !ok {
			v.Unlock()
			return fmt.Errorf(
				"popping cached VMAP dataPoint failed for %s", productID,
			)
		}

		newDataPoints.TPV.Sub(newDataPoints.TPV, droppedDataPoints.PV)
		newDataPoints.TVol.Sub(newDataPoints.TVol, droppedDataPoints.Vol)

		vwapCacheMemPool.Put(droppedDataPoints)
	}

	window.Push(newDataPoints)

	v.Unlock()

	result := types.VWAPResultMemPool.Get().(*types.VWAPResult)

	result.ProductID = productID
	result.Vwap = big.NewFloat(0)
	if newDataPoints.TVol.Cmp(bigZero) != 0 {
		result.Vwap.Quo(newDataPoints.TPV, newDataPoints.TVol)
	}

	logger.Debug("New result produced", zap.Any(productID, result))

	v.resultsQ <- result

	return nil
}

func (v *ProductsVwap) GetResultsQ() <-chan *types.VWAPResult {
	return v.resultsQ
}
