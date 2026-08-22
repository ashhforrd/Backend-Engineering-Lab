package retry

import (
	"errors"
	"time"
)

var ErrInvalidPolicy = errors.New("invalid retry policy")

type Policy struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64
}

func (p Policy) Validate() error {
	if p.MaxAttempts < 1 {
		return ErrInvalidPolicy
	}

	if p.InitialDelay <= 0 {
		return ErrInvalidPolicy
	}

	if p.MaxDelay < p.InitialDelay {
		return ErrInvalidPolicy
	}

	if p.Multiplier < 1 {
		return ErrInvalidPolicy
	}

	if p.Jitter < 0 || p.Jitter > 1 {
		return ErrInvalidPolicy
	}

	return nil
}
