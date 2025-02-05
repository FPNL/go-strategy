package pipeline

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Option func(*Pipeline)

type Pipeline struct {
	maxRetry     int
	waitStrategy func(int)
	tryFns       []func() error
	catchFns     []func() error
}

func NewPipeline(opts ...Option) *Pipeline {
	pipeline := &Pipeline{}
	for _, opt := range opts {
		opt(pipeline)
	}
	return pipeline
}

func (s *Pipeline) Then(try func() error, cancel func() error) *Pipeline {
	s.tryFns = append(s.tryFns, try)
	s.catchFns = append(s.catchFns, cancel)
	return s
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return strs[len(strs)-1]
}

func (s *Pipeline) Done() (err error) {
	for i, fn := range s.tryFns {
		if err = s.do(fn); err != nil {
			err = fmt.Errorf("execute %s err: %s", GetFunctionName(fn), err)
			if err2 := s.rollback(i); err2 != nil {
				return fmt.Errorf("%s\nrollback err: %s", err, err2)
			}

		}
	}

	return
}

func (s *Pipeline) rollback(i int) error {
	for j := i; j >= 0; j-- {
		if err := s.catchFns[j](); err != nil {
			return err
		}
	}

	return nil
}

func (s *Pipeline) do(fn func() error) error {
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
