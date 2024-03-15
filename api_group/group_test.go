package api_group

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const ErrURL = "https://error.com"

func mockAPI(url string) error {
	if url == ErrURL {
		return fmt.Errorf("error url: %s", url)
	}

	return nil
}

func mockRollbackAPI(s string) error {
	log.Println("rollback: ", s)
	return nil
}

const mockDomain = "https://example.com"

var mockAPI1 = func() error { return mockAPI(mockDomain) }
var mockRollbackAPI1 = func() error { return mockRollbackAPI(mockDomain) }

var mockAPI2 = func() error { return mockAPI(ErrURL) }
var mockRollbackAPI2 = func() error { return mockRollbackAPI(ErrURL) }

func TestApiGroup_Done(t *testing.T) {
	// given
	maxRetry := 4

	// when
	groupAPI := NewGroupAPI(
		WithMaxRetry(maxRetry),
	)

	err, rollbackErr := groupAPI.
		Then(mockAPI1, mockRollbackAPI1).
		Then(mockAPI2, mockRollbackAPI2).
		Done()

	// then
	assert.Error(t, err)
	assert.NoError(t, rollbackErr)
}

func TestWithJustWaitStrategy(t *testing.T) {
	// given
	now := time.Now()

	maxRetry := 2
	interval := 2
	givenStrategy := JustWait(time.Duration(interval) * time.Second)
	var expectWaitSeconds = float64(maxRetry * interval)

	// when
	groupAPI := NewGroupAPI(
		WithWaiStrategy(givenStrategy),
		WithMaxRetry(maxRetry),
	)

	err, rollbackErr := groupAPI.
		Then(mockAPI1, mockRollbackAPI1).
		Then(mockAPI2, mockRollbackAPI2).
		Done()

	// then
	assert.Error(t, err)
	assert.NoError(t, rollbackErr)
	assert.GreaterOrEqual(t, time.Since(now).Seconds(), expectWaitSeconds)
}

func TestWithWaitExponentStrategy(t *testing.T) {
	// given
	now := time.Now()

	maxRetry := 2
	givenStrategy := WaitExponentialBackoff(32)
	var expectWaitSeconds = float64((1+(2^(maxRetry-1)))*maxRetry) + 1

	// when
	groupAPI := NewGroupAPI(
		WithWaiStrategy(givenStrategy),
		WithMaxRetry(maxRetry),
	)

	err, rollbackErr := groupAPI.
		Then(mockAPI1, mockRollbackAPI1).
		Then(mockAPI2, mockRollbackAPI2).
		Done()

	// then
	assert.Error(t, err)
	assert.NoError(t, rollbackErr)
	assert.GreaterOrEqual(t, time.Since(now).Seconds(), expectWaitSeconds)
}
