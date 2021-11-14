package vwap_test

import (
	"math/big"
)

type TestProductDataPoint struct {
	Prices     []float64
	Volumes    []float64
	WindowSize uint16
	Product    string
}

// CalcMovWinWithF64 is a float64 VWAP moving window func calculator documenting
// the VWAP algorithm using float64 data types in GO's imperative syntax
// without the working out the big.Float algorithm. This could also be employed
// in tests to validate this service's VWAP big.float computing func.
// In that testing context care should be taken to employ price and volume
// values that produce rational VWAP results to avoid inequality results due to
// the distinct precision type characteristics between the float64 and the
// big.Float types.
func CalcMovWinWithF64(prodDataPoint TestProductDataPoint) []*big.Float {
	priceSeries, volumeSeries, windowSize := prodDataPoint.Prices, prodDataPoint.Volumes, prodDataPoint.WindowSize

	if windowSize == 0 {
		panic("Window size cannot be 0")
	}
	if len(priceSeries) != len(volumeSeries) {
		panic("The number of price and volume data points are not equal")
	}

	res := make([]*big.Float, len(priceSeries))
	var (
		iWindowSize = int(windowSize)
		v           float64
		pv          float64 // Price * Volume
		tpv         float64 // Cumulative Price*Volume or window Σ_(P*V)
		tv          float64 // Cumulative Volume or window Σ_Volume
	)

	for i, p := range priceSeries {
		v = volumeSeries[i]
		pv = p * v
		tpv += pv
		tv += v
		if i >= iWindowSize {
			// dropping the first price and volume values of the window to add
			// the new one in the total (P*V), total value.
			droppedPos := i - iWindowSize
			droppedPV := priceSeries[droppedPos] * volumeSeries[droppedPos]
			tpv -= droppedPV
			tv -= volumeSeries[droppedPos]
		}

		res[i] = big.NewFloat(0)
		if tv != 0 {
			resF64 := tpv / tv
			res[i].SetFloat64(resF64)
		}
	}

	// return the full slice to simulate the streaming generated results
	// while considering the moving window size reflected in the data
	return res
}
