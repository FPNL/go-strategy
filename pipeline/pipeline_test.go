package pipeline

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mockErr = fmt.Errorf("mock error")

func okFN() error       { return nil }
func okCancelFN() error { return nil }

func failFN() error       { return mockErr }
func failCancelFN() error { return mockErr }

func TestApiGroup_Canceled(t *testing.T) {
	// given
	maxRetry := 4

	// when
	groupAPI := NewPipeline(
		WithMaxRetry(maxRetry),
	)

	err := groupAPI.
		Then(okFN, failCancelFN).
		Then(failFN, failCancelFN).
		Done()

	// then
	assert.Error(t, err)
}

func TestApiGroup_OK(t *testing.T) {
	// given
	maxRetry := 4

	// when
	groupAPI := NewPipeline(
		WithMaxRetry(maxRetry),
	)

	err := groupAPI.
		Then(okFN, okCancelFN).
		Then(okFN, failCancelFN).
		Done()

	// then
	assert.NoError(t, err)
}

func TestWithJustWaitStrategy(t *testing.T) {
	// given
	now := time.Now()

	maxRetry := 2
	interval := 2
	givenStrategy := JustWait(time.Duration(interval) * time.Second)
	var expectWaitSeconds = float64(maxRetry * interval)

	// when
	pipeline := NewPipeline(
		WithWaitStrategy(givenStrategy),
		WithMaxRetry(maxRetry),
	)

	err := pipeline.
		Then(okFN, okCancelFN).
		Then(failFN, failCancelFN).
		Done()

	// then
	assert.Error(t, err)
	assert.GreaterOrEqual(t, time.Since(now).Seconds(), expectWaitSeconds)
}

func TestWithWaitExponentStrategy(t *testing.T) {
	// given
	now := time.Now()

	maxRetry := 2
	givenStrategy := WaitExponentialBackoff(32)
	var expectWaitSeconds = float64(3)

	// when
	pipeline := NewPipeline(
		WithWaitStrategy(givenStrategy),
		WithMaxRetry(maxRetry),
	)

	err := pipeline.
		Then(okFN, okCancelFN).
		Then(failFN, failCancelFN).
		Done()

	// then
	assert.Error(t, err)
	assert.GreaterOrEqual(t, time.Since(now).Seconds(), expectWaitSeconds)
}
