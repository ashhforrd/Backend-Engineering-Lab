package job

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestStoreMakesDefensivePayloadCopies(t *testing.T) {
	store := NewStore()
	payload := json.RawMessage(`{"value":"original"}`)

	store.Save(Job{
		ID:        "job-1",
		Type:      TypeSendEmail,
		Payload:   payload,
		Status:    StatusQueued,
		CreatedAt: time.Now(),
	})

	payload[10] = "X"[0]

	first, err := store.Get("job-1")
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if string(first.Payload) != `{"value":"original"}` {
		t.Fatalf("stored payload was mutated: %s", first.Payload)
	}

	first.Payload[10] = "Y"[0]

	second, err := store.Get("job-1")
	if err != nil {
		t.Fatalf("get job again: %v", err)
	}

	if string(second.Payload) != `{"value":"original"}` {
		t.Fatalf("returned payload mutated store: %s", second.Payload)
	}
}

func TestStoreEnforcesJobLifecycle(t *testing.T) {
	store := NewStore()
	store.Save(Job{
		ID:        "job-1",
		Type:      TypeGenerateReport,
		Status:    StatusQueued,
		CreatedAt: time.Now(),
	})

	if err := store.MarkCompleted("job-1"); !errors.Is(
		err,
		ErrInvalidTransition,
	) {
		t.Fatalf("expected invalid transition, got %v", err)
	}

	if err := store.MarkProcessing("job-1"); err != nil {
		t.Fatalf("mark processing: %v", err)
	}

	if err := store.MarkCompleted("job-1"); err != nil {
		t.Fatalf("mark completed: %v", err)
	}

	value, err := store.Get("job-1")
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if value.Status != StatusCompleted {
		t.Fatalf("expected COMPLETED, got %s", value.Status)
	}

	if value.StartedAt == nil || value.CompletedAt == nil {
		t.Fatal("expected lifecycle timestamps")
	}
}
