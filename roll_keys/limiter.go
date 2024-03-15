package rollKeys

import "golang.org/x/time/rate"

type Limiter[T any] struct {
	*rate.Limiter
	key T
}

// NewLimiter allow r executions per second
func NewLimiter[T any](obj T, r int) *Limiter[T] {
	limiter := rate.NewLimiter(rate.Limit(r), 1)

	return &Limiter[T]{
		Limiter: limiter,
		key:     obj,
	}
}
