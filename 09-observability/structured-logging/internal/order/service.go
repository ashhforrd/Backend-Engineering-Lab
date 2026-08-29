package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/logging"
)

var ErrPaymentDeclined = errors.New("payment was declined")

type CreateInput struct {
	CustomerID      string
	Email           string
	Amount          int64
	PaymentToken    string
	SimulateFailure bool
}

type Order struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerID"`
	Amount     int64  `json:"amount"`
	Status     string `json:"status"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (Order, error) {
	logger := logging.FromContext(ctx)

	logger.InfoContext(
		ctx,
		"order creation started",
		"customer_id",
		input.CustomerID,
		"amount",
		input.Amount,
	)

	if err := wait(
		ctx,
		300*time.Millisecond,
	); err != nil {
		logger.WarnContext(
			ctx,
			"order creation cancelled",
			"customer_id",
			input.CustomerID,
			"error_type",
			"CONTEXT_CANCELLED",
			"error",
			err,
		)

		return Order{}, err
	}

	if input.SimulateFailure {
		logger.ErrorContext(
			ctx,
			"payment failed",
			"customer_id",
			input.CustomerID,
			"amount",
			input.Amount,
			"error_type",
			"PAYMENT_DECLINED",
			"error",
			ErrPaymentDeclined,
		)

		return Order{}, ErrPaymentDeclined
	}

	orderID, err := generateID()
	if err != nil {
		return Order{}, err
	}

	result := Order{
		ID:         orderID,
		CustomerID: input.CustomerID,
		Amount:     input.Amount,
		Status:     "CREATED",
	}

	logger.InfoContext(
		ctx,
		"order created",
		"order_id",
		result.ID,
		"customer_id",
		result.CustomerID,
		"amount",
		result.Amount,
	)

	return result, nil
}

func generateID() (string, error) {
	value := make([]byte, 12)

	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(value), nil
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
