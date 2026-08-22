package simulator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/downstream"
)

type Handler struct {
	mu       sync.Mutex
	attempts map[string]int
}

func NewHandler() *Handler {
	return &Handler{
		attempts: make(map[string]int),
	}
}

func (h *Handler) CreatePayment(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var paymentRequest downstream.PaymentRequest

	if err := json.NewDecoder(request.Body).Decode(&paymentRequest); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	attempt := h.incrementAttempt(paymentRequest.OrderID)

	if attempt <= paymentRequest.FailuresBeforeSuccess {
		http.Error(
			writer,
			fmt.Sprintf("temporary failure on attempt: %d", attempt),
			http.StatusServiceUnavailable,
		)
		return
	}

	response := downstream.PaymentResponse{
		PaymentID: "pay-" + paymentRequest.OrderID,
		Status:    "SUCCESS",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		return
	}
}

func (h *Handler) incrementAttempt(orderID string) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.attempts[orderID]++

	return h.attempts[orderID]
}
