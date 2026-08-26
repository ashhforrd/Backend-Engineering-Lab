package job

import (
	"encoding/json"
	"time"
)

type Type string

const (
	TypeSendEmail      Type = "SEND_EMAIL"
	TypeGenerateReport Type = "GENERATE_REPORT"
)

type Status string

const (
	StatusQueued     Status = "QUEUED"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
)

type Job struct {
	ID          string          `json:"id"`
	Type        Type            `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      Status          `json:"status"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	StartedAt   *time.Time      `json:"startedAt,omitempty"`
	CompletedAt *time.Time      `json:"completedAt,omitempty"`
}
