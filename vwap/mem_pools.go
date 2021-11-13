package vwap

import (
	"sync"
)

var vwapCacheMemPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(vwapCache)
	},
}
