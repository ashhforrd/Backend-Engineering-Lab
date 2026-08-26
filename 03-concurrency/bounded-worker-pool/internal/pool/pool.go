package pool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool struct {
	config Config

	tasks   chan Task
	results chan Result

	mu      sync.RWMutex
	running bool

	workerContext context.Context
	cancelWorkers context.CancelFunc

	waitGroup sync.WaitGroup

	submitted         atomic.Int64
	activeWorkers     atomic.Int64
	maxActiveObserved atomic.Int64
	completed         atomic.Int64
	failed            atomic.Int64
}

func New(config Config) (*Pool, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Pool{
		config: config,
		tasks: make(
			chan Task,
			config.QueueCapacity,
		),
		results: make(
			chan Result,
			config.QueueCapacity,
		),
	}, nil
}

func (p *Pool) Start(
	parentContext context.Context,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return ErrAlreadyRunning
	}

	p.workerContext, p.cancelWorkers = context.WithCancel(parentContext)

	p.running = true

	for workerID := 1; workerID <= p.config.WorkerCount; workerID++ {
		p.waitGroup.Add(1)

		go p.runWorker(
			p.workerContext,
			workerID,
		)
	}

	return nil
}

func (p *Pool) Submit(
	ctx context.Context,
	task Task,
) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.running {
		return ErrNotRunning
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-p.workerContext.Done():
		return ErrNotRunning

	case p.tasks <- task:
		p.submitted.Add(1)
		return nil

	default:
		return ErrQueueFull
	}
}

func (p *Pool) Stop() {
	p.mu.Lock()

	if !p.running {
		p.mu.Unlock()
		return
	}

	p.running = false
	p.cancelWorkers()
	p.mu.Unlock()

	p.waitGroup.Wait()
	close(p.results)
}

func (p *Pool) Results() <-chan Result {
	return p.results
}
