package critical

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/lock"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	lockManager    *lock.Manager
	redisClient    *redis.Client
	instanceID     string
	acquireTimeout time.Duration
	logger         *slog.Logger
}

type Result struct {
	InstanceID string `json:"instanceID"`
	Resource   string `json:"resource"`
	Counter    int64  `json:"counter"`
}

func NewService(
	lockManager *lock.Manager,
	redisClient *redis.Client,
	instanceID string,
	acquireTimeout time.Duration,
	logger *slog.Logger,
) *Service {
	return &Service{
		lockManager:    lockManager,
		redisClient:    redisClient,
		instanceID:     instanceID,
		acquireTimeout: acquireTimeout,
		logger:         logger,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	resource string,
	workDuration time.Duration,
) (Result, error) {
	acquireContext, cancleAcquire := context.WithTimeout(
		ctx,
		s.acquireTimeout,
	)
	defer cancleAcquire()

	lease, err := s.lockManager.Acquire(
		acquireContext,
		resource,
	)
	if err != nil {
		return Result{}, err
	}

	defer s.releaseLease(lease)

	keepAliveContext, cancelKeepAlive := context.WithCancel(ctx)
	defer cancelKeepAlive()

	keepAliveErrors := lease.KeepAlive(
		keepAliveContext,
	)

	timer := time.NewTimer(workDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
	case err, open := <-keepAliveErrors:
		if !open {
			return Result{}, ctx.Err()
		}

		return Result{}, fmt.Errorf(
			"maintain lock lease: %w",
			err,
		)
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}

	counter, err := s.redisClient.Incr(
		ctx,
		"critical-counter:"+resource,
	).Result()
	if err != nil {
		return Result{}, fmt.Errorf(
			"incremental critical counter: %w", err,
		)
	}

	return Result{
		InstanceID: s.instanceID,
		Resource:   resource,
		Counter:    counter,
	}, nil
}

func (s *Service) releaseLease(
	lease *lock.Lease,
) {
	releaseContext, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	if err := lease.Release(releaseContext); err != nil && !errors.Is(err, lock.ErrNotOwner) {
		s.logger.Warn(
			"release distributed lock",
			"error",
			err,
		)
	}
}
