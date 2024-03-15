package strategy

import "sync/atomic"

func CircularPick[T any](circularIdx int64, slice []T) T {
	atomic.AddInt64(&circularIdx, 1)

	if circularIdx >= int64(len(slice)) {
		circularIdx = 0
	}

	return slice[circularIdx]
}
