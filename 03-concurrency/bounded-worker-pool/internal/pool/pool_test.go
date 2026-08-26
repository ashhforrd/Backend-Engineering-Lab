package pool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRejectsInvalidConfiguration(t *testing.T) {
	_, err := New(Config{
		WorkerCount:   0,
		QueueCapacity: 1,
	})

	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected invalid config, got %v", err)
	}
}

func TestSubmitRequiresRunningPool(t *testing.T) {
	workerPool, err := New(Config{
		WorkerCount:   1,
		QueueCapacity: 1,
	})
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	err = workerPool.Submit(
		context.Background(),
		Task{
			ID: "task-1",
			Run: func(context.Context) error {
				return nil
			},
		},
	)

	if !errors.Is(err, ErrNotRunning) {
		t.Fatalf("expected not running, got %v", err)
	}
}

func TestPoolNeverExceedsWorkerCount(t *testing.T) {
	const (
		workerCount = 3
		taskCount   = 12
	)

	workerPool, err := New(Config{
		WorkerCount:   workerCount,
		QueueCapacity: taskCount,
	})
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if err := workerPool.Start(
		context.Background(),
	); err != nil {
		t.Fatalf("start pool: %v", err)
	}

	releaseTasks := make(chan struct{})

	for taskNumber := range taskCount {
		err := workerPool.Submit(
			context.Background(),
			Task{
				ID: string(rune(taskNumber)),
				Run: func(context.Context) error {
					<-releaseTasks
					return nil
				},
			},
		)
		if err != nil {
			t.Fatalf("submit task: %v", err)
		}
	}

	waitForCondition(
		t,
		func() bool {
			return workerPool.Stats().ActiveWorkers ==
				workerCount
		},
	)

	stats := workerPool.Stats()

	if stats.MaxActiveObserved != workerCount {
		t.Fatalf(
			"expected max active %d, got %d",
			workerCount,
			stats.MaxActiveObserved,
		)
	}

	close(releaseTasks)

	for range taskCount {
		<-workerPool.Results()
	}

	workerPool.Stop()

	stats = workerPool.Stats()

	if stats.Completed != taskCount {
		t.Fatalf(
			"expected %d completed, got %d",
			taskCount,
			stats.Completed,
		)
	}
}

func TestSubmitReturnsQueueFull(t *testing.T) {
	workerPool, err := New(Config{
		WorkerCount:   1,
		QueueCapacity: 1,
	})
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if err := workerPool.Start(
		context.Background(),
	); err != nil {
		t.Fatalf("start pool: %v", err)
	}

	releaseTask := make(chan struct{})

	if err := workerPool.Submit(
		context.Background(),
		Task{
			ID: "active",
			Run: func(context.Context) error {
				<-releaseTask
				return nil
			},
		},
	); err != nil {
		t.Fatalf("submit active task: %v", err)
	}

	waitForCondition(
		t,
		func() bool {
			return workerPool.Stats().ActiveWorkers == 1
		},
	)

	if err := workerPool.Submit(
		context.Background(),
		Task{
			ID: "queued",
			Run: func(context.Context) error {
				return nil
			},
		},
	); err != nil {
		t.Fatalf("submit queued task: %v", err)
	}

	err = workerPool.Submit(
		context.Background(),
		Task{
			ID: "rejected",
			Run: func(context.Context) error {
				return nil
			},
		},
	)

	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected queue full, got %v", err)
	}

	close(releaseTask)

	for range 2 {
		<-workerPool.Results()
	}

	workerPool.Stop()
}

func TestExecuteTaskRecoversPanic(t *testing.T) {
	result := executeTask(
		context.Background(),
		1,
		Task{
			ID: "panicking-task",
			Run: func(context.Context) error {
				panic("boom")
			},
		},
	)

	if result.Err == nil {
		t.Fatal("expected panic to become an error")
	}

	if result.FinishedAt.IsZero() {
		t.Fatal("expected task finish time")
	}
}

func waitForCondition(
	t *testing.T,
	condition func() bool,
) {
	t.Helper()

	deadline := time.Now().Add(time.Second)

	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("condition was not reached")
		}

		time.Sleep(time.Millisecond)
	}
}
