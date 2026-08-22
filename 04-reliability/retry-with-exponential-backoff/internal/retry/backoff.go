package retry

import (
	"math"
	"math/rand/v2"
	"time"
)

func (p Policy) Backoff(attempt int) time.Duration {
	exponent := math.Pow(p.Multiplier, float64(attempt-1))
	delay := float64(p.InitialDelay) * exponent

	if delay > float64(p.MaxDelay) {
		delay = float64(p.MaxDelay)
	}

	if p.Jitter > 0 {
		minimum := 1 - p.Jitter
		maximum := 1 + p.Jitter
		factor := minimum + rand.Float64()*(maximum-minimum)
		delay *= factor
	}

	return time.Duration(delay)
}
