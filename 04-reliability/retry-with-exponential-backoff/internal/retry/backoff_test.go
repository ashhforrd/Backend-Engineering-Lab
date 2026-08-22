package retry

import (
	"testing"
	"time"
)

func TestPolicyBackoffGrowsExponentially(t *testing.T) {
	policy := Policy{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     500 * time.Millisecond,
		Multiplier:   2,
		Jitter:       0,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 1, want: 100 * time.Millisecond},
		{attempt: 2, want: 200 * time.Millisecond},
		{attempt: 3, want: 400 * time.Millisecond},
		{attempt: 4, want: 500 * time.Millisecond},
		{attempt: 5, want: 500 * time.Millisecond},
	}

	for _, test := range tests {
		got := policy.Backoff(test.attempt)

		if got != test.want {
			t.Errorf(
				"attempt %d: expected %s, got %s",
				test.attempt,
				test.want,
				got,
			)
		}
	}
}

func TestPolicyValidateRejectsInvalidConfiguration(t *testing.T) {
	tests := []Policy{
		{
			MaxAttempts:  0,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
			Multiplier:   2,
		},
		{
			MaxAttempts:  3,
			InitialDelay: 0,
			MaxDelay:     time.Second,
			Multiplier:   2,
		},
		{
			MaxAttempts:  3,
			InitialDelay: time.Second,
			MaxDelay:     500 * time.Millisecond,
			Multiplier:   2,
		},
		{
			MaxAttempts:  3,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
			Multiplier:   0.5,
		},
		{
			MaxAttempts:  3,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
			Multiplier:   2,
			Jitter:       1.5,
		},
	}

	for index, policy := range tests {
		if err := policy.Validate(); err == nil {
			t.Errorf(
				"case %d: expected validation error",
				index,
			)
		}
	}
}
