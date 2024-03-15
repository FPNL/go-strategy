package api_group

type Option func(*ApiGroup)

type ApiGroup struct {
	maxRetry     int
	waitStrategy func(int)
	tryFns       []func() error
	catchFns     []func() error
}

func NewGroupAPI(opts ...Option) *ApiGroup {
	apiGroup := &ApiGroup{}
	for _, opt := range opts {
		opt(apiGroup)
	}
	return apiGroup
}

func (s *ApiGroup) Then(try func() error, cancel func() error) *ApiGroup {
	s.tryFns = append(s.tryFns, try)
	s.catchFns = append(s.catchFns, cancel)
	return s
}

func (s *ApiGroup) Done() (error, error) {
	for i, fn := range s.tryFns {
		if err := s.do(fn); err != nil {
			return err, s.rollback(i)
		}
	}
	return nil, nil
}

func (s *ApiGroup) rollback(i int) error {
	for j := i; j >= 0; j-- {
		if err := s.catchFns[j](); err != nil {
			return err
		}
	}

	return nil
}

func (s *ApiGroup) do(fn func() error) error {
	for i := 0; i < s.maxRetry; i++ {
		if err := fn(); err == nil {
			return nil
		}
		if s.waitStrategy != nil {
			s.waitStrategy(i)
		}
	}

	return fn()
}
