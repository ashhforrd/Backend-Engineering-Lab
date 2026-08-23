package circuitbreaker

import "errors"

var (
	ErrCircuitOpen             = errors.New("circuit breaker is open")
	ErrTooManyHalfOpenRequests = errors.New("too many half-open requests")
)
