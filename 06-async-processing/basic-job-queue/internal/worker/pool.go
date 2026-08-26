package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/queue"
)

type Pool struct {
	workerCount int
	queue       *queue.Queue
	store       *job.Store
	dispatcher  *Dispatcher
	logger      *slog.Logger
	waitGroup   sync.WaitGroup
}

func NewPool(
	workerCount int,
	queue *queue.Queue,
	store *job.Store,
	dispatcher *Dispatcher,
	logger *slog.Logger,
) *Pool {
	return &Pool{
		workerCount: workerCount,
		queue:       queue,
		store:       store,
		dispatcher:  dispatcher,
		logger:      logger,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for workerID := 1; workerID <= p.workerCount; workerID++ {
		p.waitGroup.Add(1)

		go p.runWorker(
			ctx,
			workerID,
		)
	}
}

func (p *Pool) Wait() {
	p.waitGroup.Wait()
}

func (p *Pool) runWorker(
	ctx context.Context,
	workerID int,
) {
	defer p.waitGroup.Done()

	p.logger.InfoContext(
		ctx,
		"worker started",
		"worker_id",
		workerID,
	)

	for {
		jobID, err := p.queue.Dequeue(ctx)
		if err != nil {
			p.logger.Info(
				"worker stopped",
				"worker_id",
				workerID,
			)
			return
		}

		p.processJob(
			ctx,
			workerID,
			jobID,
		)
	}
}

func (p *Pool) processJob(
	ctx context.Context,
	workerID int,
	jobID string,
) {
	value, err := p.store.Get(jobID)
	if err != nil {
		p.logger.ErrorContext(
			ctx,
			"load job",
			"worker_id",
			workerID,
			"job_id",
			jobID,
			"error",
			err,
		)
		return
	}

	if err := p.store.MarkProcessing(jobID); err != nil {
		p.logger.ErrorContext(
			ctx,
			"mark job processing",
			"worker_id",
			workerID,
			"job_id",
			jobID,
			"error",
			err,
		)
		return
	}

	err = p.dispatcher.Dispatch(ctx, value)
	if err != nil {
		if markErr := p.store.MarkFailed(
			jobID,
			err.Error(),
		); markErr != nil {
			p.logger.ErrorContext(
				ctx,
				"mark job failed",
				"job_id",
				jobID,
				"error",
				markErr,
			)
		}

		return
	}

	if err := p.store.MarkCompleted(jobID); err != nil {
		p.logger.ErrorContext(
			ctx,
			"mark job completed",
			"job_id",
			jobID,
			"error",
			err,
		)
	}
}
