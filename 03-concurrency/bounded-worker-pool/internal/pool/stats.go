package pool

type Stats struct {
	WorkerCount       int   `json:"workerCount"`
	QueueCapacity     int   `json:"queueCapacity"`
	QueueLength       int   `json:"queueLength"`
	Submitted         int64 `json:"submitted"`
	ActiveWorkers     int64 `json:"activeWorkers"`
	MaxActiveObserved int64 `json:"maxActiveObserved"`
	Completed         int64 `json:"completed"`
	Failed            int64 `json:"failed"`
}

func (p *Pool) Stats() Stats {
	return Stats{
		WorkerCount:       p.config.WorkerCount,
		QueueCapacity:     p.config.QueueCapacity,
		QueueLength:       len(p.tasks),
		Submitted:         p.submitted.Load(),
		ActiveWorkers:     p.activeWorkers.Load(),
		MaxActiveObserved: p.maxActiveObserved.Load(),
		Completed:         p.completed.Load(),
		Failed:            p.failed.Load(),
	}
}

func (p *Pool) recordStarted() {
	active := p.activeWorkers.Add(1)

	for {
		currentMaximum := p.maxActiveObserved.Load()

		if active <= currentMaximum {
			return
		}

		if p.maxActiveObserved.CompareAndSwap(
			currentMaximum,
			active,
		) {
			return
		}
	}
}

func (p *Pool) recordFinished(err error) {
	p.activeWorkers.Add(-1)

	if err != nil {
		p.failed.Add(1)
		return
	}

	p.completed.Add(1)
}
