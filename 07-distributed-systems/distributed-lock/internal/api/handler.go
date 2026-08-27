package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/critical"
	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/lock"
)

type Handler struct {
	service *critical.Service
}

type ExecuteRequest struct {
	WorkDurationMS int `json:"workDurationMs"`
}

func NewHandler(
	service *critical.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Execute(
	writer http.ResponseWriter,
	request *http.Request,
) {
	resource := request.PathValue("resource")
	if resource == "" {
		writeError(
			writer,
			http.StatusBadRequest,
			"resource is required",
		)
		return
	}

	var input ExecuteRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if input.WorkDurationMS <= 0 || input.WorkDurationMS > 30000 {
		writeError(
			writer,
			http.StatusBadRequest,
			"workDurationMs must be between 1 and 30000",
		)
		return
	}

	result, err := h.service.Execute(
		request.Context(),
		resource,
		time.Duration(input.WorkDurationMS)*time.Millisecond,
	)
	if errors.Is(err, lock.ErrNotAcquired) {
		writeError(
			writer,
			http.StatusConflict,
			"resource is locked by another process",
		)
		return
	}

	if err != nil {
		writeError(
			writer,
			http.StatusInternalServerError,
			"critical section failed",
		)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		result,
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
