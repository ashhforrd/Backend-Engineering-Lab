package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/payment"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/retry"
)

type Handler struct {
	paymentService *payment.Service
}

func NewHandler(paymentService *payment.Service) *Handler {
	return &Handler{
		paymentService: paymentService,
	}
}

func (h *Handler) CreatePayment(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var paymentRequest downstream.PaymentRequest

	if err := json.NewDecoder(request.Body).Decode(&paymentRequest); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid request body")
		return
	}

	if paymentRequest.OrderID == "" {
		writeError(writer, http.StatusBadRequest, "orderId is required")
		return
	}

	if paymentRequest.Amount <= 0 {
		writeError(writer, http.StatusBadRequest, "amount must be greater than zero")
		return
	}

	if paymentRequest.FailuresBeforeSuccess < 0 {
		writeError(
			writer,
			http.StatusBadRequest,
			"failuresBeforeSuccess cannot be negative",
		)
		return
	}

	result, err := h.paymentService.Create(
		request.Context(),
		paymentRequest,
	)
	if err != nil {
		if errors.Is(err, retry.ErrAttemptsExhausted) {
			writeError(
				writer,
				http.StatusBadGateway,
				"downstream service remained unavailable",
			)
			return
		}

		writeError(
			writer,
			http.StatusBadGateway,
			"downstream request failed",
		)
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func writeError(
	writer http.ResponseWriter,
	statusCode int,
	message string,
) {
	writeJSON(
		writer,
		statusCode,
		map[string]string{"error": message},
	)
}

func writeJSON(
	writer http.ResponseWriter,
	statusCode int,
	value any,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(value); err != nil {
		return
	}
}
