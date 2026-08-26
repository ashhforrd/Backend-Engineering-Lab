package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/03-concurrency/bounded-worker-pool/internal/pool"
)

type Handler struct {
	pool      *pool.Pool
	idCounter atomic.Int64
}

type SubmitRequest struct {
	TaskCount  int `json:"taskCount"`
	DurationMS int `json:"durationMs"`
	FailEvery  int `json:"failEvery"`
}

type SubmitResponse struct {
	Accepted int `json:"accepted"`
	Rejected int `json:"rejected"`
}

func NewHandler(workerPool *pool.Pool) *Handler {
	return &Handler{
		pool: workerPool,
	}
}

func (h *Handler) SubmitTasks(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input SubmitRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if input.TaskCount <= 0 ||
		input.TaskCount > 1000 ||
		input.DurationMS <= 0 ||
		input.DurationMS > 30000 ||
		input.FailEvery < 0 {

		writeError(
			writer,
			http.StatusBadRequest,
			"invalid task configuration",
		)
		return
	}

	response := SubmitResponse{}

	for taskNumber := 1; taskNumber <= input.TaskCount; taskNumber++ {

		taskID := strconv.FormatInt(
			h.idCounter.Add(1),
			10,
		)

		taskNumber := taskNumber

		err := h.pool.Submit(
			request.Context(),
			pool.Task{
				ID: taskID,
				Run: func(ctx context.Context) error {
					if err := wait(
						ctx,
						time.Duration(input.DurationMS)*
							time.Millisecond,
					); err != nil {
						return err
					}

					if input.FailEvery > 0 &&
						taskNumber%input.FailEvery == 0 {
						return errors.New(
							"simulated task failure",
						)
					}

					return nil
				},
			},
		)
		if err != nil {
			response.Rejected++
			continue
		}

		response.Accepted++
	}

	writeJSON(
		writer,
		http.StatusAccepted,
		response,
	)
}

func (h *Handler) GetStats(
	writer http.ResponseWriter,
	request *http.Request,
) {
	writeJSON(
		writer,
		http.StatusOK,
		h.pool.Stats(),
	)
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

func writeError(
	writer http.ResponseWriter,
	statusCode int,
	message string,
) {
	writeJSON(
		writer,
		statusCode,
		map[string]string{
			"error": message,
		},
	)
}

func writeJSON(
	writer http.ResponseWriter,
	statusCode int,
	value any,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)
	writer.WriteHeader(statusCode)

	_ = json.NewEncoder(writer).Encode(value)
}
