package retry

import (
	"context"
	"errors"
	"time"
)

var ErrAttemptsExhausted = errors.New("retry attempts exhausted")

type Operation[T any] func(ctx context.Context) (T, error)

type ShouldRetry func(err error) bool

func Do[T any](
	ctx context.Context,
	policy Policy,
	operation Operation[T],
	shouldRetry ShouldRetry,
) (T, error) {
	var zero T

	if err := policy.Validate(); err != nil {
		return zero, err
	}

	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt += 1 {
		if err := ctx.Err(); err != nil {
			return zero, err
		}

		result, err := operation(ctx)
		if err == nil {
			return result, nil
		}

		lastErr = err

		if !shouldRetry(err) {
			return zero, err
		}

		if attempt == policy.MaxAttempts {
			break
		}

		delay := policy.Backoff(attempt)

		if err := wait(ctx, delay); err != nil {
			return zero, err
		}
	}

	return zero, errors.Join(ErrAttemptsExhausted, lastErr)
}

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
