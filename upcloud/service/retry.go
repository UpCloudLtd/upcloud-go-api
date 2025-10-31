package service

import (
	"context"
	"errors"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const maxRetriesOnError = 3

type RetryConfig struct {
	interval time.Duration
	// Inverse the should retry logic. By default, operation is retried until operation returns a value. If inverse is set to true, operation is retried while operation returns a value. This should be used, for example, for waiting until resource is deleted.
	inverse bool
}

func fillDefaults(c *RetryConfig) *RetryConfig {
	if c == nil {
		c = &RetryConfig{}
	}

	if c.interval.Milliseconds() == 0 {
		c.interval = time.Second * 5
	}

	return c
}

func shouldRetryOnError(err error, retryOnErrorCount int) bool {
	if retryOnErrorCount >= maxRetriesOnError {
		return false
	}

	var ucErr *upcloud.Problem
	if errors.As(err, &ucErr) && ucErr.Status >= 500 {
		return true
	}

	return false
}

func Retry[T any](ctx context.Context, operation func(int, context.Context) (*T, error), config *RetryConfig) (*T, error) {
	config = fillDefaults(config)

	ticker := time.NewTicker(config.interval)
	defer ticker.Stop()

	retryOnErrorCount := 0
	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			value, err := operation(i, ctx)

			if err != nil {
				if shouldRetryOnError(err, retryOnErrorCount) {
					retryOnErrorCount++
					continue
				}

				return value, err
			} else {
				retryOnErrorCount = 0
			}

			if !config.inverse && value != nil {
				return value, nil
			}

			if config.inverse && value == nil {
				return nil, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
