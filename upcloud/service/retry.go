package service

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
)

type waitConfig struct {
	interval time.Duration
	// Inverse the should retry logic. By default, operation is retried until operation returns a value. If inverse is set to true, operation is retried while operation returns a value. This should be used, for example, for waiting until resource is deleted.
	inverse bool
}

func fillDefaults(c *waitConfig) *waitConfig {
	if c == nil {
		c = &waitConfig{}
	}

	if c.interval.Milliseconds() == 0 {
		c.interval = time.Second * 5
	}

	return c
}

func withWaitFragment(ctx context.Context, i int) context.Context {
	return client.WithFragment(ctx, fmt.Sprintf("wait:i=%d", i))
}

func wait[T any](ctx context.Context, operation func(int, context.Context) (*T, error), config *waitConfig) (*T, error) {
	config = fillDefaults(config)

	ticker := time.NewTicker(config.interval)
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			c := withWaitFragment(ctx, i)
			value, err := operation(i, c)
			if err != nil {
				return value, err
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
