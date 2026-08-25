package product

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
)

type repositoryStub struct {
	findByID func(context.Context, int64) (Product, error)
	update   func(context.Context, Product) (Product, error)
}

func (r *repositoryStub) FindByID(
	ctx context.Context,
	productID int64,
) (Product, error) {
	return r.findByID(ctx, productID)
}

func (r *repositoryStub) Update(
	ctx context.Context,
	value Product,
) (Product, error) {
	return r.update(ctx, value)
}

type cacheStub struct {
	get    func(context.Context, int64) (Product, error)
	set    func(context.Context, Product, time.Duration) error
	delete func(context.Context, int64) error
}

func (c *cacheStub) Get(
	ctx context.Context,
	productID int64,
) (Product, error) {
	return c.get(ctx, productID)
}

func (c *cacheStub) Set(
	ctx context.Context,
	value Product,
	ttl time.Duration,
) error {
	return c.set(ctx, value, ttl)
}

func (c *cacheStub) Delete(
	ctx context.Context,
	productID int64,
) error {
	return c.delete(ctx, productID)
}

func testLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}

func TestGetReturnsCachedProductWithoutCallingRepository(
	t *testing.T,
) {
	cached := Product{
		ID:    1,
		Name:  "Cached Keyboard",
		Price: 1200000,
		Stock: 10,
	}

	repositoryCalled := false

	repository := &repositoryStub{
		findByID: func(
			context.Context,
			int64,
		) (Product, error) {
			repositoryCalled = true
			return Product{}, nil
		},
		update: func(
			context.Context,
			Product,
		) (Product, error) {
			return Product{}, nil
		},
	}

	cache := &cacheStub{
		get: func(
			context.Context,
			int64,
		) (Product, error) {
			return cached, nil
		},
		set: func(
			context.Context,
			Product,
			time.Duration,
		) error {
			return nil
		},
		delete: func(
			context.Context,
			int64,
		) error {
			return nil
		},
	}

	service := NewService(
		repository,
		cache,
		30*time.Second,
		testLogger(),
	)

	result, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get product: %v", err)
	}

	if result.Source != SourceCache {
		t.Fatalf("expected CACHE, got %s", result.Source)
	}

	if result.Product != cached {
		t.Fatalf(
			"expected %#v, got %#v",
			cached,
			result.Product,
		)
	}

	if repositoryCalled {
		t.Fatal("repository must not be called on cache hit")
	}
}

func TestGetLoadsAndCachesProductOnCacheMiss(
	t *testing.T,
) {
	databaseProduct := Product{
		ID:    1,
		Name:  "Database Keyboard",
		Price: 1200000,
		Stock: 10,
	}

	var cachedProduct Product
	var cachedTTL time.Duration

	repository := &repositoryStub{
		findByID: func(
			context.Context,
			int64,
		) (Product, error) {
			return databaseProduct, nil
		},
		update: func(
			context.Context,
			Product,
		) (Product, error) {
			return Product{}, nil
		},
	}

	cache := &cacheStub{
		get: func(
			context.Context,
			int64,
		) (Product, error) {
			return Product{}, ErrCacheMiss
		},
		set: func(
			_ context.Context,
			value Product,
			ttl time.Duration,
		) error {
			cachedProduct = value
			cachedTTL = ttl
			return nil
		},
		delete: func(
			context.Context,
			int64,
		) error {
			return nil
		},
	}

	service := NewService(
		repository,
		cache,
		30*time.Second,
		testLogger(),
	)

	result, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get product: %v", err)
	}

	if result.Source != SourceDatabase {
		t.Fatalf("expected DATABASE, got %s", result.Source)
	}

	if cachedProduct != databaseProduct {
		t.Fatalf(
			"expected cached product %#v, got %#v",
			databaseProduct,
			cachedProduct,
		)
	}

	if cachedTTL != 30*time.Second {
		t.Fatalf("expected 30s TTL, got %s", cachedTTL)
	}
}

func TestGetFallsBackToRepositoryWhenCacheFails(
	t *testing.T,
) {
	databaseProduct := Product{
		ID:   1,
		Name: "Available from database",
	}

	repository := &repositoryStub{
		findByID: func(
			context.Context,
			int64,
		) (Product, error) {
			return databaseProduct, nil
		},
		update: func(
			context.Context,
			Product,
		) (Product, error) {
			return Product{}, nil
		},
	}

	cache := &cacheStub{
		get: func(
			context.Context,
			int64,
		) (Product, error) {
			return Product{}, errors.New("redis unavailable")
		},
		set: func(
			context.Context,
			Product,
			time.Duration,
		) error {
			return errors.New("redis unavailable")
		},
		delete: func(
			context.Context,
			int64,
		) error {
			return nil
		},
	}

	service := NewService(
		repository,
		cache,
		30*time.Second,
		testLogger(),
	)

	result, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get product: %v", err)
	}

	if result.Source != SourceDatabase {
		t.Fatalf("expected DATABASE, got %s", result.Source)
	}

	if result.Product != databaseProduct {
		t.Fatalf(
			"expected %#v, got %#v",
			databaseProduct,
			result.Product,
		)
	}
}

func TestUpdateInvalidatesCachedProduct(t *testing.T) {
	input := Product{
		ID:    1,
		Name:  "Keyboard Pro",
		Price: 1500000,
		Stock: 8,
	}

	deletedProductID := int64(0)

	repository := &repositoryStub{
		findByID: func(
			context.Context,
			int64,
		) (Product, error) {
			return Product{}, nil
		},
		update: func(
			_ context.Context,
			value Product,
		) (Product, error) {
			return value, nil
		},
	}

	cache := &cacheStub{
		get: func(
			context.Context,
			int64,
		) (Product, error) {
			return Product{}, ErrCacheMiss
		},
		set: func(
			context.Context,
			Product,
			time.Duration,
		) error {
			return nil
		},
		delete: func(
			_ context.Context,
			productID int64,
		) error {
			deletedProductID = productID
			return nil
		},
	}

	service := NewService(
		repository,
		cache,
		30*time.Second,
		testLogger(),
	)

	result, err := service.Update(
		context.Background(),
		input,
	)
	if err != nil {
		t.Fatalf("update product: %v", err)
	}

	if result != input {
		t.Fatalf("expected %#v, got %#v", input, result)
	}

	if deletedProductID != input.ID {
		t.Fatalf(
			"expected deleted product %d, got %d",
			input.ID,
			deletedProductID,
		)
	}
}
