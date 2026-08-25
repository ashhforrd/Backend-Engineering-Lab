package circuitbreaker

import (
	"sync"
	"time"
)

type Breaker struct {
	mu sync.Mutex

	config Config
	state  State

	consecutiveFailures uint32
	openedAt            time.Time

	halfOpenInFlight  uint32
	halfOpenSuccesses uint32

	generation uint64
	now        func() time.Time
}

func New(config Config) (*Breaker, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Breaker{
		config: config,
		state:  StateClosed,
		now:    time.Now,
	}, nil
}

func (b *Breaker) transitionToOpen() {
	b.state = StateOpen
	b.openedAt = b.now()

	b.consecutiveFailures = 0
	b.halfOpenInFlight = 0
	b.halfOpenSuccesses = 0

	b.generation++
}

func (b *Breaker) transitionToHalfOpen() {
	b.state = StateHalfOpen

	b.consecutiveFailures = 0
	b.halfOpenInFlight = 0
	b.halfOpenSuccesses = 0

	b.generation++
}

func (b *Breaker) transitionToClosed() {
	b.state = StateClosed
	b.openedAt = time.Time{}

	b.consecutiveFailures = 0
	b.halfOpenInFlight = 0
	b.halfOpenSuccesses = 0

	b.generation++
}
