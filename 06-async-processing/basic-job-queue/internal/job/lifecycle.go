package job

import (
	"errors"
	"fmt"
	"time"
)

var ErrInvalidTransition = errors.New(
	"invalid job status transition",
)

func (s *Store) MarkProcessing(
	jobID string,
) error {
	return s.transition(
		jobID,
		StatusQueued,
		StatusProcessing,
		"",
	)
}

func (s *Store) MarkCompleted(
	jobID string,
) error {
	return s.transition(
		jobID,
		StatusProcessing,
		StatusCompleted,
		"",
	)
}

func (s *Store) MarkFailed(
	jobID string,
	message string,
) error {
	return s.transition(
		jobID,
		StatusProcessing,
		StatusFailed,
		message,
	)
}

func (s *Store) transition(
	jobID string,
	expected Status,
	next Status,
	errorMesssage string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.jobs[jobID]
	if !exists {
		return ErrNotFound
	}

	if value.Status != expected {
		return fmt.Errorf(
			"%w: expected=%s actual=%s next=%s",
			ErrInvalidTransition,
			expected,
			value.Status,
			next,
		)
	}

	now := time.Now().UTC()
	value.Status = next

	switch next {
	case StatusProcessing:
		value.StartedAt = &now
		value.Error = ""
	case StatusCompleted, StatusFailed:
		value.CompletedAt = &now
		value.Error = errorMesssage
	}

	s.jobs[jobID] = value

	return nil
}
