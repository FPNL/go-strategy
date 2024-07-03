package rollKeys

import (
	"context"
	"errors"
)

type RotationalSlice[T any] struct {
	s            []*Limiter[T]
	pickStrategy func([]*Limiter[T]) *Limiter[T]
}

func NewRotationalSlice[T any](slice []T, r int, options ...func(*RotationalSlice[T])) (*RotationalSlice[T], error) {
	if len(slice) == 0 {
		return nil, errors.New("empty slice")
	}
	if r <= 0 {
		return nil, errors.New("rate must be positive")
	}

	var ls = make([]*Limiter[T], len(slice))

	for i, obj := range slice {
		ls[i] = NewLimiter(obj, r)
	}

	var count int64 = 0

	c := &RotationalSlice[T]{
		s:            ls,
		pickStrategy: CircularPick[*Limiter[T]](&count),
	}

	for _, option := range options {
		option(c)
	}

	return c, nil
}

func (c *RotationalSlice[T]) Get(ctx context.Context) (T, error) {
	for _, l := range c.s {
		if l.Allow() {
			return l.key, nil
		}
	}

	l := c.pick()
	err := l.Wait(ctx)
	if err != nil {
		return l.key, err
	}

	return l.key, nil
}

func (c *RotationalSlice[T]) pick() *Limiter[T] {
	return c.pickStrategy(c.s)
}
