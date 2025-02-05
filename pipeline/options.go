package pipeline

func WithMaxRetry(i int) Option {
	return func(s *Pipeline) {
		s.maxRetry = i
	}
}

func WithWaitStrategy(fn func(int)) Option {
	return func(s *Pipeline) {
		s.waitStrategy = fn
	}
}
