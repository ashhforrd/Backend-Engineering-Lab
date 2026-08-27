package lock

import (
	"context"
	"time"
)

func (l *Lease) KeepAlive(
	ctx context.Context,
) <-chan error {
	result := make(chan error, 1)

	go func() {
		defer close(result)

		interval := l.manager.config.TTL / 3
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				if err := l.Extend(ctx); err != nil {
					result <- err
					return
				}
			}
		}
	}()

	return result
}
