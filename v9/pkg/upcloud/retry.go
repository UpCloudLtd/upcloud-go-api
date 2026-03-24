package upcloud

import (
	"context"
	"errors"
	"time"
)

const (
	defaultPollInterval = 5 * time.Second
	maxRetriesOnError   = 3
)

func Retryable(err error) error {
	if err == nil {
		return nil
	}
	return &retryableError{err: err}
}

type retryableError struct{ err error }

func (e *retryableError) Error() string { return e.err.Error() }
func (e *retryableError) Unwrap() error { return e.err }

func retry[T any](ctx context.Context, op func(int, context.Context) (*T, error)) (*T, error) {
	return poll(ctx, op, func(v *T) bool { return v != nil })
}

func retryUntilNil[T any](ctx context.Context, op func(int, context.Context) (*T, error)) error {
	_, err := poll(ctx, op, func(v *T) bool { return v == nil })
	return err
}

func poll[T any](ctx context.Context, op func(int, context.Context) (*T, error), done func(*T) bool) (*T, error) {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	retryCount := 0
	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			value, err := op(i, ctx)
			if err != nil {
				if isRetryableError(err) && retryCount < maxRetriesOnError {
					retryCount++
					continue
				}
				return value, err
			}

			retryCount = 0

			if done(value) {
				return value, nil
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func isRetryableError(err error) bool {
	var re *retryableError
	return errors.As(err, &re)
}
