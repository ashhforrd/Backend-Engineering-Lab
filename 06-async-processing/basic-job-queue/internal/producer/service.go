package producer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/queue"
)

var ErrInvalidJobType = errors.New(
	"invalid job type",
)

type Service struct {
	store *job.Store
	queue *queue.Queue
}

func NewService(
	store *job.Store,
	queue *queue.Queue,
) *Service {
	return &Service{
		store: store,
		queue: queue,
	}
}

func (s *Service) Create(
	ctx context.Context,
	jobType job.Type,
	payload json.RawMessage,
) (job.Job, error) {
	if !isSupportedType(jobType) {
		return job.Job{}, ErrInvalidJobType
	}

	jobID, err := generateID()
	if err != nil {
		return job.Job{}, err
	}

	value := job.Job{
		ID:        jobID,
		Type:      jobType,
		Payload:   payload,
		Status:    job.StatusQueued,
		CreatedAt: time.Now().UTC(),
	}

	s.store.Save(value)

	if err := s.queue.Enqueue(
		ctx,
		value.ID,
	); err != nil {
		s.store.Delete(value.ID)
		return job.Job{}, err
	}

	return value, nil
}

func isSupportedType(jobType job.Type) bool {
	switch jobType {
	case job.TypeSendEmail, job.TypeGenerateReport:
		return true

	default:
		return false
	}
}

func generateID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
