package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(
	pool *pgxpool.Pool,
) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) FindByID(
	ctx context.Context,
	productID int64,
) (product.Product, error) {
	var result product.Product

	err := r.pool.QueryRow(
		ctx,
		`
			SELECT id, name, price, stock, updated_at
			FROM products
			WHERE id = $1
		`,
		productID,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Price,
		&result.Stock,
		&result.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return product.Product{}, product.ErrNotFound
	}

	if err != nil {
		return product.Product{}, fmt.Errorf(
			"find product by id: %w",
			err,
		)
	}

	return result, nil
}

func (r *Repository) Update(
	ctx context.Context,
	input product.Product,
) (product.Product, error) {
	var result product.Product

	err := r.pool.QueryRow(
		ctx,
		`
			UPDATE products
			SET
				name = $2,
				price = $3,
				stock = $4,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING id, name, price, stock, updated_at
		`,
		input.ID,
		input.Name,
		input.Price,
		input.Stock,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Price,
		&result.Stock,
		&result.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return product.Product{}, product.ErrNotFound
	}

	if err != nil {
		return product.Product{}, fmt.Errorf(
			"update product: %w",
			err,
		)
	}

	return result, nil
}
