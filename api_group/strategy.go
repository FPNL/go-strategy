package api_group

import (
	"time"

	"github.com/FPNL/go-helper/strategy"
)

func JustWait(duration time.Duration) func(int) {
	return func(int) {
		time.Sleep(duration)
	}
}

func WaitExponentialBackoff(maxWaitSecond float64) func(int) {
	return func(n int) {
		time.Sleep(strategy.ExponentialBackoff(n, maxWaitSecond))
	}
}
