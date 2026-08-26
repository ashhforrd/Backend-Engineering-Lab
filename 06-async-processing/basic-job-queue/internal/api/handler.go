package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/producer"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/queue"
)

type Handler struct {
	producer *producer.Service
	store    *job.Store
	queue    *queue.Queue
}

type CreateJobRequest struct {
	Type    job.Type        `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewHandler(
	producer *producer.Service,
	store *job.Store,
	queue *queue.Queue,
) *Handler {
	return &Handler{
		producer: producer,
		store:    store,
		queue:    queue,
	}
}

func (h *Handler) CreateJob(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input CreateJobRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if len(input.Payload) == 0 {
		writeError(
			writer,
			http.StatusBadRequest,
			"payload is required",
		)
		return
	}

	value, err := h.producer.Create(
		request.Context(),
		input.Type,
		input.Payload,
	)
	if errors.Is(err, producer.ErrInvalidJobType) {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid job type",
		)
		return
	}

	if errors.Is(err, queue.ErrFull) {
		writeError(
			writer,
			http.StatusBadRequest,
			"job queue is full",
		)
		return
	}

	if err != nil {
		writeError(
			writer,
			http.StatusInternalServerError,
			"create job failed",
		)
		return
	}

	writer.Header().Set(
		"Location",
		"/api/jobs"+value.ID,
	)

	writeJSON(
		writer,
		http.StatusAccepted,
		value,
	)
}

func (h *Handler) GetJob(
	writer http.ResponseWriter,
	request *http.Request,
) {
	value, err := h.store.Get(
		request.PathValue("id"),
	)
	if errors.Is(err, job.ErrNotFound) {
		writeError(
			writer,
			http.StatusNotFound,
			"job not found",
		)
		return
	}

	if err != nil {
		writeError(
			writer,
			http.StatusInternalServerError,
			"get job failed",
		)
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

func (h *Handler) GetQueueStatus(
	writer http.ResponseWriter,
	request *http.Request,
) {
	writeJSON(
		writer,
		http.StatusOK,
		map[string]int{
			"length":   h.queue.Length(),
			"capacity": h.queue.Capacity(),
		},
	)
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
