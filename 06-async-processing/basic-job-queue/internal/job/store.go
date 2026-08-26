package job

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("job not found")

type Store struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewStore() *Store {
	return &Store{
		jobs: make(map[string]Job),
	}
}

func (s *Store) Save(value Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value.Payload = clonePayload(value.Payload)
	s.jobs[value.ID] = value
}

func (s *Store) Get(jobID string) (Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.jobs[jobID]
	if !exists {
		return Job{}, ErrNotFound
	}

	value.Payload = clonePayload(value.Payload)

	return value, nil
}

func clonePayload(payload []byte) []byte {
	cloned := make([]byte, len(payload))
	copy(cloned, payload)

	return cloned
}

func (s *Store) Delete(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.jobs, jobID)
}
