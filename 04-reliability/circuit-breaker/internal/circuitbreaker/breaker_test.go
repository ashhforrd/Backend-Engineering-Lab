package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBreakerTransitionsThroughLifecycle(t *testing.T) {
	currentTime := time.Date(
		2026,
		time.August,
		25,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	breaker, err := New(Config{
		FailureThreshold:     2,
		OpenTimeout:          5 * time.Second,
		HalfOpenMaxRequests:  1,
		HalfOpenSuccessLimit: 1,
	})
	if err != nil {
		t.Fatalf("create breaker: %v", err)
	}

	breaker.now = func() time.Time {
		return currentTime
	}

	downstreamErr := errors.New("downstream unavailable")
	calls := 0

	failingOperation := func(
		ctx context.Context,
	) (string, error) {
		calls++
		return "", downstreamErr
	}

	isFailure := func(err error) bool {
		return true
	}

	for range 2 {
		_, _ = Execute(
			context.Background(),
			breaker,
			failingOperation,
			isFailure,
		)
	}

	if breaker.Snapshot().State != StateOpen {
		t.Fatalf(
			"expected OPEN, got %s",
			breaker.Snapshot().State,
		)
	}

	_, err = Execute(
		context.Background(),
		breaker,
		failingOperation,
		isFailure,
	)

	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf(
			"expected circuit open error, got %v",
			err,
		)
	}

	if calls != 2 {
		t.Fatalf(
			"expected downstream to be called twice, got %d",
			calls,
		)
	}

	currentTime = currentTime.Add(5 * time.Second)

	result, err := Execute(
		context.Background(),
		breaker,
		func(ctx context.Context) (string, error) {
			calls++
			return "success", nil
		},
		isFailure,
	)
	if err != nil {
		t.Fatalf("expected successful probe, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected success, got %q", result)
	}

	if breaker.Snapshot().State != StateClosed {
		t.Fatalf(
			"expected CLOSED, got %s",
			breaker.Snapshot().State,
		)
	}
}
