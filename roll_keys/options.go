package rollKeys

func WithPickStrategy[T any](pickStrategy func([]*Limiter[T]) *Limiter[T]) func(*RotationalSlice[T]) {
	return func(c *RotationalSlice[T]) {
		c.pickStrategy = pickStrategy
	}
}
