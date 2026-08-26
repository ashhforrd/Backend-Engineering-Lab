package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
)

var ErrRequestedFailure = errors.New(
	"job requested to fail",
)

type Handlers struct {
	logger *slog.Logger
}

type SendEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Fail    bool   `json:"fail"`
}

type GenerateReportPayload struct {
	ReportName string `json:"reportName"`
	Fail       bool   `json:"fail"`
}

func NewHandlers(
	logger *slog.Logger,
) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) SendEmail(
	ctx context.Context,
	value job.Job,
) error {
	var payload SendEmailPayload

	if err := json.Unmarshal(
		value.Payload,
		&payload,
	); err != nil {
		return fmt.Errorf(
			"decode email payload: %w",
			err,
		)
	}

	if payload.To == "" || payload.Subject == "" {
		return errors.New(
			"email recipient and subject are required",
		)
	}

	if err := wait(ctx, 500*time.Millisecond); err != nil {
		return err
	}

	if payload.Fail {
		return ErrRequestedFailure
	}

	h.logger.InfoContext(
		ctx,
		"email sent",
		"job_id",
		value.ID,
		"to",
		payload.To,
	)

	return nil
}

func (h *Handlers) GenerateReport(
	ctx context.Context,
	value job.Job,
) error {
	var payload GenerateReportPayload

	if err := json.Unmarshal(
		value.Payload,
		&payload,
	); err != nil {
		return fmt.Errorf(
			"decode report payload: %w",
			err,
		)
	}

	if payload.ReportName == "" {
		return errors.New(
			"report name is required",
		)
	}

	if err := wait(ctx, time.Second); err != nil {
		return err
	}

	if payload.Fail {
		return ErrRequestedFailure
	}

	h.logger.InfoContext(
		ctx,
		"report generated",
		"job_id",
		value.ID,
		"report_name",
		payload.ReportName,
	)

	return nil
}

func wait(
	ctx context.Context,
	duration time.Duration,
) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
