package rollKeys

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestRotationalSlice_Get(t *testing.T) {
	now := time.Now()

	// given
	APIKeys := []string{"api-key-1", "api-key-2"}
	requestTimes := 50
	rate := 2

	expectDuration := requestTimes / (len(APIKeys) * rate)

	// when
	keys, err := NewRotationalSlice(APIKeys, rate)
	require.NoError(t, err)

	eg := errgroup.Group{}

	for i := 0; i < requestTimes; i++ {
		eg.Go(func() error {
			key, err := keys.Get(context.TODO())
			if err != nil {
				return err
			}

			assert.Contains(t, APIKeys, key)

			return nil
		})
	}
	err = eg.Wait()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectDuration, int(time.Since(now).Seconds()))
}
