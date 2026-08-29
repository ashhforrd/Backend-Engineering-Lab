package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/order"
)

type Handler struct {
	orderService *order.Service
}

type CreateOrderRequest struct {
	CustomerID      string `json:"customerId"`
	Email           string `json:"email"`
	Amount          int64  `json:"amount"`
	PaymentToken    string `json:"paymentToken"`
	SimulateFailure bool   `json:"simulateFailure"`
}

func NewHandler(
	orderService *order.Service,
) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h *Handler) CreateOrder(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input CreateOrderRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if input.CustomerID == "" ||
		input.Email == "" ||
		input.Amount <= 0 ||
		input.PaymentToken == "" {

		writeError(
			writer,
			http.StatusBadRequest,
			"invalid order data",
		)
		return
	}

	result, err := h.orderService.Create(
		request.Context(),
		order.CreateInput{
			CustomerID:      input.CustomerID,
			Email:           input.Email,
			Amount:          input.Amount,
			PaymentToken:    input.PaymentToken,
			SimulateFailure: input.SimulateFailure,
		},
	)
	if errors.Is(err, order.ErrPaymentDeclined) {
		writeError(
			writer,
			http.StatusUnprocessableEntity,
			"payment was declined",
		)
		return
	}

	if err != nil {
		writeError(
			writer,
			http.StatusInternalServerError,
			"create order failed",
		)
		return
	}

	writeJSON(
		writer,
		http.StatusCreated,
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
