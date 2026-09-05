package ratelimit

import (
	"testing"
	"time"
)

type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	return c.now
}

func (c *fakeClock) Advance(duration time.Duration) {
	c.now = c.now.Add(duration)
}

func TestLimiterConsumesAndRejectsQuota(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1_000, 0)}
	limiter, err := NewLimiter(Config{Limit: 2, Window: time.Minute}, clock)
	if err != nil {
		t.Fatalf("create limiter: %v", err)
	}

	first := limiter.Allow("client-1")
	second := limiter.Allow("client-1")
	third := limiter.Allow("client-1")

	if !first.Allowed || first.Remaining != 1 {
		t.Fatalf("unexpected first decision: %+v", first)
	}

	if !second.Allowed || second.Remaining != 0 {
		t.Fatalf("unexpected second decision: %+v", second)
	}

	if third.Allowed || third.Remaining != 0 || third.RetryAfter != time.Minute {
		t.Fatalf("unexpected rejected decision: %+v", third)
	}
}

func TestLimiterTracksClientsIndependently(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1_000, 0)}
	limiter, err := NewLimiter(Config{Limit: 1, Window: time.Minute}, clock)
	if err != nil {
		t.Fatalf("create limiter: %v", err)
	}

	if !limiter.Allow("client-1").Allowed {
		t.Fatal("expected client-1 to be allowed")
	}

	if !limiter.Allow("client-2").Allowed {
		t.Fatal("expected client-2 to have independent quota")
	}
}

func TestLimiterResetsExpiredWindow(t *testing.T) {
	clock := &fakeClock{now: time.Unix(1_000, 0)}
	limiter, err := NewLimiter(Config{Limit: 1, Window: time.Minute}, clock)
	if err != nil {
		t.Fatalf("create limiter: %v", err)
	}

	limiter.Allow("client-1")
	if limiter.Allow("client-1").Allowed {
		t.Fatal("expected exhausted quota to be rejected")
	}

	clock.Advance(time.Minute)
	decision := limiter.Allow("client-1")

	if !decision.Allowed || decision.Remaining != 0 {
		t.Fatalf("expected a fresh window, got %+v", decision)
	}
}

func TestNewLimiterValidatesConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   error
	}{
		{name: "zero limit", config: Config{Window: time.Minute}, want: ErrInvalidLimit},
		{name: "zero window", config: Config{Limit: 1}, want: ErrInvalidWindow},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewLimiter(test.config, &fakeClock{})
			if err != test.want {
				t.Fatalf("expected %v, got %v", test.want, err)
			}
		})
	}
}
