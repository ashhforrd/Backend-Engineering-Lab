package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
)

var ErrUnsupportedJobType = errors.New("unsupported job type")

type Handler func(
	ctx context.Context,
	value job.Job,
) error

type Dispatcher struct {
	handlers map[job.Type]Handler
}

func NewDispatcher(
	handlers map[job.Type]Handler,
) *Dispatcher {
	return &Dispatcher{
		handlers: handlers,
	}
}

func (d *Dispatcher) Dispatch(
	ctx context.Context,
	value job.Job,
) error {
	handler, exists := d.handlers[value.Type]
	if !exists {
		return fmt.Errorf(
			"%w: %s",
			ErrUnsupportedJobType,
			value.Type,
		)
	}

	return handler(ctx, value)
}
