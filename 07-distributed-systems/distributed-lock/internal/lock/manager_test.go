package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestOnlyOneLeaseCanOwnResource(t *testing.T) {
	manager, _, cleanup := newTestManager(
		t,
		time.Second,
		5*time.Millisecond,
	)
	defer cleanup()

	firstLease, err := manager.Acquire(
		context.Background(),
		"daily-report",
	)
	if err != nil {
		t.Fatalf("acquire first lease: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		20*time.Millisecond,
	)
	defer cancel()

	_, err = manager.Acquire(ctx, "daily-report")
	if !errors.Is(err, ErrNotAcquired) {
		t.Fatalf("expected not acquired, got %v", err)
	}

	if err := firstLease.Release(
		context.Background(),
	); err != nil {
		t.Fatalf("release first lease: %v", err)
	}

	secondLease, err := manager.Acquire(
		context.Background(),
		"daily-report",
	)
	if err != nil {
		t.Fatalf("acquire second lease: %v", err)
	}

	if err := secondLease.Release(
		context.Background(),
	); err != nil {
		t.Fatalf("release second lease: %v", err)
	}
}

func TestExpiredOwnerCannotReleaseNewOwnersLock(
	t *testing.T,
) {
	const ttl = time.Second

	manager, server, cleanup := newTestManager(
		t,
		ttl,
		5*time.Millisecond,
	)
	defer cleanup()

	expiredLease, err := manager.Acquire(
		context.Background(),
		"daily-report",
	)
	if err != nil {
		t.Fatalf("acquire expiring lease: %v", err)
	}

	server.FastForward(ttl)

	currentLease, err := manager.Acquire(
		context.Background(),
		"daily-report",
	)
	if err != nil {
		t.Fatalf("acquire current lease: %v", err)
	}

	err = expiredLease.Release(context.Background())
	if !errors.Is(err, ErrNotOwner) {
		t.Fatalf("expected not owner, got %v", err)
	}

	value, err := manager.client.Get(
		context.Background(),
		currentLease.key,
	).Result()
	if err != nil {
		t.Fatalf("read current lock: %v", err)
	}

	if value != currentLease.token {
		t.Fatal("expired owner removed the current owners lock")
	}
}

func TestExtendRenewsLeaseTTL(t *testing.T) {
	const ttl = time.Second

	manager, server, cleanup := newTestManager(
		t,
		ttl,
		5*time.Millisecond,
	)
	defer cleanup()

	lease, err := manager.Acquire(
		context.Background(),
		"daily-report",
	)
	if err != nil {
		t.Fatalf("acquire lease: %v", err)
	}

	server.FastForward(750 * time.Millisecond)

	if err := lease.Extend(
		context.Background(),
	); err != nil {
		t.Fatalf("extend lease: %v", err)
	}

	server.FastForward(750 * time.Millisecond)

	exists, err := manager.client.Exists(
		context.Background(),
		lease.key,
	).Result()
	if err != nil {
		t.Fatalf("check lock: %v", err)
	}

	if exists != 1 {
		t.Fatal("expected extended lease to remain active")
	}
}

func TestNewManagerRejectsInvalidConfig(t *testing.T) {
	client := redis.NewClient(
		&redis.Options{
			Addr: "unused",
		},
	)
	defer client.Close()

	_, err := NewManager(
		client,
		Config{},
	)

	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected invalid config, got %v", err)
	}
}

func newTestManager(
	t *testing.T,
	ttl time.Duration,
	retryInterval time.Duration,
) (*Manager, *miniredis.Miniredis, func()) {
	t.Helper()

	server := miniredis.RunT(t)

	client := redis.NewClient(
		&redis.Options{
			Addr: server.Addr(),
		},
	)

	manager, err := NewManager(
		client,
		Config{
			TTL:           ttl,
			RetryInterval: retryInterval,
		},
	)
	if err != nil {
		t.Fatalf("create manager: %v", err)
	}

	cleanup := func() {
		_ = client.Close()
		server.Close()
	}

	return manager, server, cleanup
}
