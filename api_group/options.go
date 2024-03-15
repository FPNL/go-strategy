package api_group

func WithMaxRetry(i int) Option {
	return func(s *ApiGroup) {
		s.maxRetry = i
	}
}

func WithWaiStrategy(fn func(int)) Option {
	return func(s *ApiGroup) {
		s.waitStrategy = fn
	}
}
