package pool

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrNilTask = errors.New(
	"task function is nil",
)

func (p *Pool) runWorker(
	ctx context.Context,
	workerID int,
) {
	defer p.waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case task := <-p.tasks:
			p.recordStarted()

			result := executeTask(
				ctx,
				workerID,
				task,
			)

			p.recordFinished(result.Err)

			select {
			case p.results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func executeTask(
	ctx context.Context,
	workerID int,
	task Task,
) (result Result) {
	result.TaskID = task.ID
	result.WorkerID = workerID
	result.StartedAt = time.Now().UTC()

	defer func() {
		result.FinishedAt = time.Now().UTC()

		if recovered := recover(); recovered != nil {
			result.Err = fmt.Errorf(
				"task panicked: %v",
				recovered,
			)
		}
	}()

	if task.Run == nil {
		result.Err = ErrNilTask
		return result
	}

	result.Err = task.Run(ctx)

	return result
}
