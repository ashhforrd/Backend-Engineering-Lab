package pool

import (
	"context"
	"time"
)

type TaskFunc func(ctx context.Context) error

type Task struct {
	ID  string
	Run TaskFunc
}

type Result struct {
	TaskID     string
	WorkerID   int
	StartedAt  time.Time
	FinishedAt time.Time
	Err        error
}

func (r Result) Duration() time.Duration {
	return r.FinishedAt.Sub(r.StartedAt)
}
