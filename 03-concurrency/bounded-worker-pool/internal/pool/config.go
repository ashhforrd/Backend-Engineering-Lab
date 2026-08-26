package pool

import "errors"

var (
	ErrInvalidConfig = errors.New(
		"invalid worker pool config",
	)

	ErrNotRunning = errors.New(
		"worker pool is not running",
	)

	ErrQueueFull = errors.New(
		"worker pool queue is full",
	)

	ErrAlreadyRunning = errors.New(
		"worker pool iss already running",
	)
)

type Config struct {
	WorkerCount   int
	QueueCapacity int
}

func (c Config) Validate() error {
	if c.WorkerCount <= 0 {
		return ErrInvalidConfig
	}

	if c.QueueCapacity <= 0 {
		return ErrInvalidConfig
	}

	return nil
}
