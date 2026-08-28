package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/08-infrastructure/graceful-http-server/internal/lifecycle"
)

type Handler struct {
	state      *lifecycle.State
	middleware *lifecycle.Middleware
}

func NewHandler(
	state *lifecycle.State,
	middleware *lifecycle.Middleware,
) *Handler {
	return &Handler{
		state:      state,
		middleware: middleware,
	}
}

func (h *Handler) Liveness(
	writer http.ResponseWriter,
	requesst *http.Request,
) {
	writeJSON(
		writer,
		http.StatusOK,
		map[string]any{
			"status": "ALIVE",
		},
	)
}

func (h *Handler) Readiness(
	writer http.ResponseWriter,
	request *http.Request,
) {
	statusCode := http.StatusOK
	status := "READY"

	if !h.state.IsReady() {
		statusCode = http.StatusServiceUnavailable
		status = "NOT_READY"
	}

	writeJSON(
		writer,
		statusCode,
		map[string]any{
			"statuss":        status,
			"activeRequests": h.middleware.ActiveRequests(),
		},
	)
}

func (h *Handler) SlowWork(
	writer http.ResponseWriter,
	request *http.Request,
) {
	durationMilliseconds, err := strconv.Atoi(
		request.URL.Query().Get("durationMs"),
	)
	if err != nil ||
		durationMilliseconds <= 0 ||
		durationMilliseconds > 30000 {

		writeError(
			writer,
			http.StatusBadRequest,
			"durationMs must be between 1 and 30000",
		)
		return
	}

	duration := time.Duration(
		durationMilliseconds,
	) * time.Millisecond

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		writeJSON(
			writer,
			http.StatusOK,
			map[string]any{
				"status":     "COMPLETED",
				"durationMs": durationMilliseconds,
			},
		)

	case <-request.Context().Done():
		return
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
		map[string]any{
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
