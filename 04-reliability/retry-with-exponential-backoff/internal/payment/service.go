package payment

import (
	"context"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/retry"
)

type Service struct {
	client *downstream.Client
	policy retry.Policy
}

type Result struct {
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
	Attempts  int    `json"attempts"`
}

func NewService(
	client *downstream.Client,
	policy retry.Policy,
) *Service {
	return &Service{
		client: client,
		policy: policy,
	}
}

func (s *Service) Create(
	ctx context.Context,
	request downstream.PaymentRequest,
) (Result, error) {
	attempts := 0

	response, err := retry.Do(
		ctx,
		s.policy,
		func(ctx context.Context) (downstream.PaymentResponse, error) {
			attempts++

			return s.client.CreatePayment(ctx, request)
		},
		downstream.IsRetryable,
	)
	if err != nil {
		return Result{}, err
	}

	return Result{
		PaymentID: response.PaymentID,
		Status:    response.Status,
		Attempts:  attempts,
	}, nil
}
