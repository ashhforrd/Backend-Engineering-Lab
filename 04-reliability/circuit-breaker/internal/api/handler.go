package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/circuitbreaker"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/product"
)

type Handler struct {
	productService *product.Service
	breaker        *circuitbreaker.Breaker
}

func NewHandler(
	productService *product.Service,
	breaker *circuitbreaker.Breaker,
) *Handler {
	return &Handler{
		productService: productService,
		breaker:        breaker,
	}
}

func (h *Handler) GetProduct(
	writer http.ResponseWriter,
	request *http.Request,
) {
	productID := request.PathValue("id")

	result, err := h.productService.Get(
		request.Context(),
		productID,
	)
	if err != nil {
		h.writeRequestError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (h *Handler) GetCircuitStatus(
	writer http.ResponseWriter,
	request *http.Request,
) {
	writeJSON(
		writer,
		http.StatusOK,
		h.breaker.Snapshot(),
	)
}

func (h *Handler) writeRequestError(
	writer http.ResponseWriter,
	err error,
) {
	if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
		writeError(
			writer,
			http.StatusServiceUnavailable,
			"circuit breaker is open",
		)
		return
	}

	if errors.Is(
		err,
		circuitbreaker.ErrTooManyHalfOpenRequests,
	) {
		writeError(
			writer,
			http.StatusServiceUnavailable,
			"half-open request limit reached",
		)
		return
	}

	var downstreamErr *downstream.Error

	if errors.As(err, &downstreamErr) {
		writeError(
			writer,
			http.StatusBadGateway,
			downstreamErr.Message,
		)
		return
	}

	writeError(
		writer,
		http.StatusInternalServerError,
		"internal server error",
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