package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry_noInverse(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	value, err := retry(ctx, func(i int, _ context.Context) (*string, error) {
		if i < 3 {
			return nil, nil
		}

		value := "ready"
		return &value, nil
	}, &retryConfig{interval: time.Millisecond * 125})

	assert.NoError(t, err)
	assert.Equal(t, "ready", *value)
}

func TestRetry_timeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	value, err := retry(ctx, func(i int, _ context.Context) (*string, error) {
		return nil, nil
	}, &retryConfig{interval: time.Millisecond * 125})

	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestRetry_inverse(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	value, err := retry(ctx, func(i int, _ context.Context) (*int, error) {
		if i < 3 {
			return &i, nil
		}
		return nil, nil
	}, &retryConfig{interval: time.Millisecond * 125, inverse: true})

	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestRetry_noInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *retryConfig
	}{
		{name: "no interval", config: &retryConfig{inverse: false}},
		{name: "nil config", config: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 125*time.Millisecond)
			defer cancel()

			value, err := retry(ctx, func(i int, _ context.Context) (*int, error) {
				return &i, nil
			}, test.config)

			assert.NoError(t, err)
			require.NotNil(t, value)
			assert.Equal(t, 0, *value)
		})
	}
}

func problemsThenValue(status, count int) func(i int, _ context.Context) (*string, error) {
	return func(i int, _ context.Context) (*string, error) {
		if i < count {
			return nil, &upcloud.Problem{Status: status}
		}

		value := "ready"
		return &value, nil
	}
}

func TestRetry_retry500(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		status      int
		count       int
		expectValue bool
		expectError bool
	}{
		{status: 500, count: 3, expectValue: true, expectError: false},
		{status: 500, count: 4, expectValue: false, expectError: true},
		{status: 499, count: 1, expectValue: false, expectError: true},
	}

	for _, test := range testcases {
		name := fmt.Sprintf("HTTP %d x %d", test.status, test.count)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			value, err := retry(ctx, problemsThenValue(test.status, test.count), &retryConfig{interval: time.Millisecond * 125})

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if test.expectValue {
				assert.Equal(t, "ready", *value)
			} else {
				assert.Nil(t, value)
			}
		})
	}
}

func TestRetry_deadline(t *testing.T) {
	t.Parallel()

	t.Run("Before interval", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second*2)
		defer cancel()

		wantValue := 0
		gotValue, err := retry(ctx, func(i int, ctx context.Context) (*int, error) {
			return &i, ctx.Err()
		}, &retryConfig{interval: time.Second * 10})

		require.NoError(t, err)
		require.NotNil(t, gotValue)
		require.Equal(t, wantValue, *gotValue)
	})

	t.Run("After interval", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second*2)
		defer cancel()

		c := 0
		gotValue, err := retry(ctx, func(i int, ctx context.Context) (*int, error) {
			c++
			return nil, nil
		}, &retryConfig{interval: time.Second})

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, gotValue)
		require.Greater(t, c, 1)
	})

	t.Run("Reached", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), -time.Second)
		defer cancel()

		wantValue := 0
		gotValue, err := retry(ctx, func(i int, ctx context.Context) (*int, error) {
			return &i, ctx.Err()
		}, &retryConfig{interval: time.Second * 10})

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NotNil(t, gotValue)
		require.Equal(t, wantValue, *gotValue)
	})
}
