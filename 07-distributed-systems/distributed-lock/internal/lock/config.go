package lock

import (
	"errors"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid distributed lock config")
	ErrNotAcquired   = errors.New("distributed lock was not acquired")
	ErrNotOwner      = errors.New("distributed lock is not owned by this process")
	ErrLockLost      = errors.New("distributed lock ownership was lost")
)

type Config struct {
	TTL           time.Duration
	RetryInterval time.Duration
}

func (c Config) Validate() error {
	if c.TTL <= 0 {
		return ErrInvalidConfig
	}

	if c.RetryInterval <= 0 {
		return ErrInvalidConfig
	}

	if c.RetryInterval >= c.TTL {
		return ErrInvalidConfig
	}

	return nil
}
