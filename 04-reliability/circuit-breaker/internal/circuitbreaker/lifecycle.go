package circuitbreaker

type ticket struct {
	generation uint64
	state      State
}

func (b *Breaker) beforeRequest() (ticket, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen &&
		b.now().Sub(b.openedAt) >= b.config.OpenTimeout {
		b.transitionToHalfOpen()
	}

	switch b.state {
	case StateOpen:
		return ticket{}, ErrCircuitOpen

	case StateHalfOpen:
		if b.halfOpenInFlight >= b.config.HalfOpenMaxRequests {
			return ticket{}, ErrTooManyHalfOpenRequests
		}

		b.halfOpenInFlight++
	}

	return ticket{
		generation: b.generation,
		state:      b.state,
	}, nil
}

func (b *Breaker) afterSuccess(requestTicket ticket) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if requestTicket.generation != b.generation {
		return
	}

	switch b.state {
	case StateClosed:
		b.consecutiveFailures = 0

	case StateHalfOpen:
		b.halfOpenInFlight--
		b.halfOpenSuccesses++

		if b.halfOpenSuccesses >= b.config.HalfOpenSuccessLimit {
			b.transitionToClosed()
		}
	}
}

func (b *Breaker) afterFailure(requestTicket ticket) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if requestTicket.generation != b.generation {
		return
	}

	switch b.state {
	case StateClosed:
		b.consecutiveFailures++

		if b.consecutiveFailures >= b.config.FailureThreshold {
			b.transitionToOpen()
		}

	case StateHalfOpen:
		b.halfOpenInFlight--
		b.transitionToOpen()
	}
}
