package vwap

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/blewater/zh/log"
	"github.com/blewater/zh/types"
	"go.uber.org/zap"
)

// Cache of cumulative data points to calculate a fresh VWAP result for each
// new product price, volume pairs.
type vwapCache struct {
	TPV  *big.Float
	TVol *big.Float
	PV   *big.Float
	Vol  *big.Float
}

// Memory pool of vwapCache objects
var vwapCacheMemPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(vwapCache)
	},
}

// ProductsVwap is the container for calculating the queued results.
type ProductsVwap struct {
	windowSize uint16
	vwapCache  sync.Map
	resultsQ   types.ResultsQ
}

var bigZero = big.NewFloat(0)

func New(productIDs []string, windowSize uint16) *ProductsVwap {
	prodVwap := &ProductsVwap{
		vwapCache:  sync.Map{},
		resultsQ:   make(types.ResultsQ, windowSize),
		windowSize: windowSize,
	}
	for _, p := range productIDs {
		// Allocate capacity upfront
		prodVwap.vwapCache.Store(p, NewWindowQueue(windowSize))
	}

	return prodVwap
}

// ProduceVwap is the service VWAP computing func employing big.Float data types.
// See CalcMovWinWithF64() for the exact same algorithm in simpler terms.
func (v *ProductsVwap) ProduceVwap(ctx context.Context, productID string, price, volume *big.Float) error {
	logger := log.FromContext(ctx)
	// nolint:errcheck
	defer logger.Sync()

	newDataPoints := memPoolGet()

	newDataPoints.PV.Mul(price, volume)
	newDataPoints.Vol.Set(volume)
	newDataPoints.TPV.Set(newDataPoints.PV)
	newDataPoints.TVol.Set(newDataPoints.Vol)

	i, ok := v.vwapCache.Load(productID)
	if !ok {
		return fmt.Errorf(
			"product ID %s not in the VWAP map of product ids", productID,
		)
	}
	window, ok := i.(*WindowQueue)
	if !ok {
		return fmt.Errorf(
			"failed to access the VWAP window slice for %s", productID,
		)
	}

	//---------------- Start a product's VWAP computation using shared memory containers
	window.Lock()
	if window.len > 0 {
		prevDataPoints, ok := window.PeekLast()
		if !ok {
			return fmt.Errorf(
				"could not access cached data set %d, %s", window.len,
				productID,
			)
		}

		// Add previous sums
		newDataPoints.TPV.Add(newDataPoints.PV, prevDataPoints.TPV)
		newDataPoints.TVol.Add(newDataPoints.Vol, prevDataPoints.TVol)
	}

	// drop window data point to make room for the new
	var droppedDataPoints *vwapCache
	if window.len == v.windowSize {
		droppedDataPoints, ok = window.Pop()
		if !ok {
			return fmt.Errorf(
				"popping cached VMAP dataPoint failed for %s", productID,
			)
		}

		newDataPoints.TPV.Sub(newDataPoints.TPV, droppedDataPoints.PV)
		newDataPoints.TVol.Sub(newDataPoints.TVol, droppedDataPoints.Vol)
	}

	window.Push(newDataPoints)
	window.Unlock()
	//---------------- End of product's VWAP computation using shared memory containers

	recycleToPool(droppedDataPoints)

	result := types.VWAPResultMemPool.Get().(*types.VWAPResult)

	result.ProductID = productID
	result.Vwap = big.NewFloat(0)
	if newDataPoints.TVol.Cmp(bigZero) != 0 {
		result.Vwap.Quo(newDataPoints.TPV, newDataPoints.TVol)
	}

	logger.Debug("New result produced", zap.Object(productID, result))

	v.resultsQ <- result

	return nil
}

func memPoolGet() *vwapCache {
	newDataPoints := vwapCacheMemPool.Get().(*vwapCache)

	newDataPoints.TPV = types.BigFloatMemPool.Get().(*big.Float)
	newDataPoints.TVol = types.BigFloatMemPool.Get().(*big.Float)
	newDataPoints.PV = types.BigFloatMemPool.Get().(*big.Float)
	newDataPoints.Vol = types.BigFloatMemPool.Get().(*big.Float)
	return newDataPoints
}

func recycleToPool(droppedDataPoints *vwapCache) {
	if droppedDataPoints != nil {
		types.BigFloatMemPool.Put(droppedDataPoints.TPV)
		types.BigFloatMemPool.Put(droppedDataPoints.TVol)
		types.BigFloatMemPool.Put(droppedDataPoints.PV)
		types.BigFloatMemPool.Put(droppedDataPoints.Vol)
		vwapCacheMemPool.Put(droppedDataPoints)
	}
}

func (v *ProductsVwap) GetResultsQ() <-chan *types.VWAPResult {
	return v.resultsQ
}
