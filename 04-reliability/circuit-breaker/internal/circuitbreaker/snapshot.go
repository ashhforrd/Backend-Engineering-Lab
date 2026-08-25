package circuitbreaker

import "time"

type Snapshot struct {
	State               State      `json:"state"`
	ConsecutiveFailures uint32     `json:"consecutiveFailures"`
	HalfOpenInFlight    uint32     `json:"halfOpenInFlight"`
	HalfOpenSuccesses   uint32     `json:"halfOpenSuccesses"`
	OpenedAt            *time.Time `json:"openedAt,omitempty"`
	Generation          uint64     `json:"generation"`
}

func (b *Breaker) Snapshot() Snapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	var openedAt *time.Time

	if !b.openedAt.IsZero() {
		value := b.openedAt
		openedAt = &value
	}

	return Snapshot{
		State:               b.state,
		ConsecutiveFailures: b.consecutiveFailures,
		HalfOpenInFlight:    b.halfOpenInFlight,
		HalfOpenSuccesses:   b.halfOpenSuccesses,
		OpenedAt:            openedAt,
		Generation:          b.generation,
	}
}