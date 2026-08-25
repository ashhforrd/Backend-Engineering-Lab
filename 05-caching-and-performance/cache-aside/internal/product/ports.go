package product

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrCacheMiss = errors.New("product cache miss")
)

type Repository interface {
	FindByID(
		ctx context.Context,
		productID int64,
	) (Product, error)

	Update(
		ctx context.Context,
		product Product,
	) (Product, error)
}

type Cache interface {
	Get(
		ctx context.Context,
		productID int64,
	) (Product, error)

	Set(
		ctx context.Context,
		product Product,
		ttl time.Duration,
	) error

	Delete(
		ctx context.Context,
		productID int64,
	) error
}
