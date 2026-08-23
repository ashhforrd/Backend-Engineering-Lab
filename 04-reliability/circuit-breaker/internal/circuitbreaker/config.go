package circuitbreaker

import (
	"errors"
	"time"
)

var ErrInvalidConfig = errors.New("invalid circuit breaker config")

type Config struct {
	FailureThreshold     uint32
	OpenTimeout          time.Duration
	HalfOpenMaxRequests  uint32
	HalfOpenSuccessLimit uint32
}

func (c Config) Validate() error {
	if c.FailureThreshold == 0 {
		return ErrInvalidConfig
	}

	if c.OpenTimeout <= 0 {
		return ErrInvalidConfig
	}

	if c.HalfOpenMaxRequests == 0 {
		return ErrInvalidConfig
	}

	if c.HalfOpenSuccessLimit == 0 {
		return ErrInvalidConfig
	}

	if c.HalfOpenSuccessLimit > c.HalfOpenMaxRequests {
		return ErrInvalidConfig
	}

	return nil
}
