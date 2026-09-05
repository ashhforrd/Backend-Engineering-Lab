package ratelimit

import (
	"errors"
	"time"
)

var (
	ErrInvalidLimit  = errors.New("limit must be greater than zero")
	ErrInvalidWindow = errors.New("window must be greater than zero")
)

type Config struct {
	Limit  int
	Window time.Duration
}

func (c *Config) Validate() error {
	if c.Limit <= 0 {
		return ErrInvalidLimit
	}

	if c.Window <= 0 {
		return ErrInvalidWindow
	}

	return nil
}
