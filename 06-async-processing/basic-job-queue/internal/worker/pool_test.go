package worker

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/queue"
)

func TestWorkerProcessesQueuedJob(t *testing.T) {
	store := job.NewStore()
	jobQueue := queue.New(1)

	value := job.Job{
		ID:        "job-1",
		Type:      job.TypeSendEmail,
		Payload:   json.RawMessage(`{"to":"student@example.com"}`),
		Status:    job.StatusQueued,
		CreatedAt: time.Now(),
	}

	store.Save(value)

	if err := jobQueue.Enqueue(
		context.Background(),
		value.ID,
	); err != nil {
		t.Fatalf("enqueue job: %v", err)
	}

	dispatcher := NewDispatcher(
		map[job.Type]Handler{
			job.TypeSendEmail: func(
				context.Context,
				job.Job,
			) error {
				return nil
			},
		},
	)

	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	pool := NewPool(
		1,
		jobQueue,
		store,
		dispatcher,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)

	deadline := time.Now().Add(time.Second)

	for {
		result, err := store.Get(value.ID)
		if err != nil {
			t.Fatalf("get job: %v", err)
		}

		if result.Status == job.StatusCompleted {
			break
		}

		if time.Now().After(deadline) {
			t.Fatalf(
				"job did not complete, status=%s",
				result.Status,
			)
		}

		time.Sleep(time.Millisecond)
	}

	cancel()
	pool.Wait()
}
