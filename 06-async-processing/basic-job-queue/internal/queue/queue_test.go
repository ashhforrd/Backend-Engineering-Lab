package queue

import (
	"context"
	"errors"
	"testing"
)

func TestQueueRejectsJobWhenCapacityIsFull(t *testing.T) {
	value := New(1)

	if err := value.Enqueue(
		context.Background(),
		"job-1",
	); err != nil {
		t.Fatalf("enqueue first job: %v", err)
	}

	if err := value.Enqueue(
		context.Background(),
		"job-2",
	); !errors.Is(err, ErrFull) {
		t.Fatalf("expected queue full, got %v", err)
	}

	jobID, err := value.Dequeue(context.Background())
	if err != nil {
		t.Fatalf("dequeue job: %v", err)
	}

	if jobID != "job-1" {
		t.Fatalf("expected job-1, got %s", jobID)
	}
}

func TestDequeueStopsWhenContextIsCancelled(t *testing.T) {
	value := New(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := value.Dequeue(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}
