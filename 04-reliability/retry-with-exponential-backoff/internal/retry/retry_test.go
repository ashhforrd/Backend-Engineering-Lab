package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func testPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     4 * time.Millisecond,
		Multiplier:   2,
		Jitter:       0,
	}
}

func TestDoRetriesUntilOperationSucceeds(t *testing.T) {
	attempts := 0
	temporaryErr := errors.New("temporary failure")

	result, err := Do(
		context.Background(),
		testPolicy(),
		func(ctx context.Context) (string, error) {
			attempts++

			if attempts < 3 {
				return "", temporaryErr
			}

			return "success", nil
		},
		func(err error) bool {
			return errors.Is(err, temporaryErr)
		},
	)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result success, got %q", result)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestDoDoesNotRetryPermanentError(t *testing.T) {
	attempts := 0
	permanentErr := errors.New("permanent failure")

	_, err := Do(
		context.Background(),
		testPolicy(),
		func(ctx context.Context) (string, error) {
			attempts++
			return "", permanentErr
		},
		func(err error) bool {
			return false
		},
	)

	if !errors.Is(err, permanentErr) {
		t.Fatalf("expected permanent error, got: %v", err)
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestDoReturnsErrorWhenAttemptsAreExhausted(t *testing.T) {
	attempts := 0
	temporaryErr := errors.New("temporary failure")

	_, err := Do(
		context.Background(),
		testPolicy(),
		func(ctx context.Context) (string, error) {
			attempts++
			return "", temporaryErr
		},
		func(err error) bool {
			return true
		},
	)

	if !errors.Is(err, ErrAttemptsExhausted) {
		t.Fatalf("expected attempts exhausted error, got: %v", err)
	}

	if !errors.Is(err, temporaryErr) {
		t.Fatalf("expected original error to be preserved, got: %v", err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestDoStopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	attempts := 0

	_, err := Do(
		ctx,
		testPolicy(),
		func(ctx context.Context) (string, error) {
			attempts++
			return "", errors.New("should not run")
		},
		func(err error) bool {
			return true
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got: %v", err)
	}

	if attempts != 0 {
		t.Fatalf("expected 0 attempts, got %d", attempts)
	}
}
