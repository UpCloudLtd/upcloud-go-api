package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_noInverse(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	value, err := wait(ctx, func(i int, _ context.Context) (*string, error) {
		if i < 3 {
			return nil, nil
		}

		value := "ready"
		return &value, nil
	}, &waitConfig{interval: time.Millisecond * 125})

	assert.NoError(t, err)
	assert.Equal(t, "ready", *value)
}

func TestRetry_timeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	value, err := wait(ctx, func(i int, _ context.Context) (*string, error) {
		return nil, nil
	}, &waitConfig{interval: time.Millisecond * 125})

	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestRetry_inverse(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	value, err := wait(ctx, func(i int, _ context.Context) (*int, error) {
		if i < 3 {
			return &i, nil
		}
		return nil, nil
	}, &waitConfig{interval: time.Millisecond * 125, inverse: true})

	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestRetry_noInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *waitConfig
	}{
		{name: "no interval", config: &waitConfig{inverse: false}},
		{name: "nil config", config: nil},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 125*time.Millisecond)
			defer cancel()

			value, err := wait(ctx, func(i int, _ context.Context) (*int, error) {
				return &i, nil
			}, test.config)

			assert.Error(t, err)
			assert.Nil(t, value)
		})
	}
}
