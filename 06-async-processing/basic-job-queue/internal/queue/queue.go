package queue

import (
	"context"
	"errors"
)

var ErrFull = errors.New("job queue is full")

type Queue struct {
	jobs chan string
}

func New(capacity int) *Queue {
	return &Queue{
		jobs: make(chan string, capacity),
	}
}

func (q *Queue) Enqueue(
	ctx context.Context,
	jobID string,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case q.jobs <- jobID:
		return nil

	default:
		return ErrFull
	}
}

func (q *Queue) Dequeue(
	ctx context.Context,
) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()

	case jobID := <-q.jobs:
		return jobID, nil
	}
}

func (q *Queue) Length() int {
	return len(q.jobs)
}

func (q *Queue) Capacity() int {
	return cap(q.jobs)
}
