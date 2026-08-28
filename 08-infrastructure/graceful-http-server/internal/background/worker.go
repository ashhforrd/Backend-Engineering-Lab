package background

import (
	"context"
	"log/slog"
	"time"
)

type Worker struct {
	logger   *slog.Logger
	interval time.Duration
}

func NewWorker(
	logger *slog.Logger,
	interval time.Duration,
) *Worker {
	return &Worker{
		logger:   logger,
		interval: interval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("background worker started")

	for {
		select {
		case <-ticker.C:
			w.logger.Info(
				"background heartbeat",
			)
		case <-ctx.Done():
			w.logger.Info(
				"background worker stopped",
			)
			return
		}
	}
}
