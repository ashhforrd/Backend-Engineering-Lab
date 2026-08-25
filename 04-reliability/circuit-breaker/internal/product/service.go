package product

import (
	"context"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/circuitbreaker"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/downstream"
)

type Service struct {
	client *downstream.Client
	breaker *circuitbreaker.Breaker
}

func NewService(
	client *downstream.Client,
	breaker *circuitbreaker.Breaker,
) *Service {
	return &Service{
		client: client,
		breaker: breaker,
	}
}

func (s *Service) Get(
	ctx context.Context,
	productID string,
) (downstream.Product, error) {
	return circuitbreaker.Execute(
		ctx,
		s.breaker,
		func(ctx context.Context) (
			downstream.Product,
			error,
		) {
			return s.client.GetProduct(
				ctx,
				productID,
			)
		},
		downstream.IsFailure,
	)
}