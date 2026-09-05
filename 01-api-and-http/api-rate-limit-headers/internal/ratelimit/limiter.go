package ratelimit

import (
	"sync"
	"time"
)

type clientWindow struct {
	used    int
	resetAt time.Time
}

type Decision struct {
	Allowed    bool
	Limit      int
	Remaining  int
	ResetAt    time.Time
	RetryAfter time.Duration
}

type Limiter struct {
	config  Config
	clock   Clock
	mu      sync.Mutex
	clients map[string]*clientWindow
}

func NewLimiter(
	config Config,
	clock Clock,
) (*Limiter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Limiter{
		config:  config,
		clock:   clock,
		clients: make(map[string]*clientWindow),
	}, nil
}

func (l *Limiter) Allow(clientID string) Decision {
	now := l.clock.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	currentWindow, exists := l.clients[clientID]

	if !exists || !now.Before(currentWindow.resetAt) {
		currentWindow = &clientWindow{
			resetAt: now.Add(l.config.Window),
		}

		l.clients[clientID] = currentWindow
	}

	if currentWindow.used >= l.config.Limit {
		return Decision{
			Allowed:    false,
			Limit:      l.config.Limit,
			Remaining:  0,
			ResetAt:    currentWindow.resetAt,
			RetryAfter: currentWindow.resetAt.Sub(now),
		}
	}

	currentWindow.used++

	return Decision{
		Allowed:   true,
		Limit:     l.config.Limit,
		Remaining: l.config.Limit - currentWindow.used,
		ResetAt:   currentWindow.resetAt,
	}
}
