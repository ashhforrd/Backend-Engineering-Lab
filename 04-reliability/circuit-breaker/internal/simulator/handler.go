package simulator

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/downstream"
)

type Handler struct {
	available atomic.Bool
}

type AvailabilityRequest struct {
	Available bool `json:"available"`
}

func NewHandler() *Handler {
	handler := &Handler{}
	handler.available.Store(true)

	return handler
}

func (h *Handler) GetProduct(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if !h.available.Load() {
		http.Error(
			writer,
			"dependency unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	productID := request.PathValue("id")

	response := downstream.Product{
		ID:    productID,
		Name:  "Mechanical Keyboard",
		Price: 1200000,
	}

	writeJSON(writer, http.StatusOK, response)
}

func (h *Handler) SetAvailability(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input AvailabilityRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(
			writer,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	h.available.Store(input.Available)

	writeJSON(
		writer,
		http.StatusOK,
		map[string]bool{
			"available": input.Available,
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