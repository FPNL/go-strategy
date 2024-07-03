package rollKeys

import (
	"github.com/FPNL/go-strategy/strategy"
)

func CircularPick[T any](circularIdx *int64) func(slice []T) T {
	return func(slice []T) T {
		return strategy.CircularPick[T](circularIdx, slice)
	}
}
