package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/product"
	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *goredis.Client
}

func NewCache(
	client *goredis.Client,
) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(
	ctx context.Context,
	productID int64,
) (product.Product, error) {
	var result product.Product

	value, err := c.client.Get(
		ctx,
		cacheKey(productID),
	).Bytes()
	if errors.Is(err, goredis.Nil) {
		return result, product.ErrCacheMiss
	}

	if err != nil {
		return result, fmt.Errorf(
			"get product cache: %w",
			err,
		)
	}

	if err := json.Unmarshal(value, &result); err != nil {
		return product.Product{}, fmt.Errorf(
			"decode product cache: %w",
			err,
		)
	}

	return result, nil
}

func (c *Cache) Set(
	ctx context.Context,
	value product.Product,
	ttl time.Duration,
) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(
			"encode product cache: %w",
			err,
		)
	}

	if err := c.client.Set(
		ctx,
		cacheKey(value.ID),
		encoded,
		ttl,
	).Err(); err != nil {
		return fmt.Errorf(
			"set product cache: %w",
			err,
		)
	}

	return nil
}

func (c *Cache) Delete(
	ctx context.Context,
	productID int64,
) error {
	if err := c.client.Del(
		ctx,
		cacheKey(productID),
	).Err(); err != nil {
		return fmt.Errorf(
			"delete product cache: %w",
			err,
		)
	}

	return nil
}

func cacheKey(productID int64) string {
	return fmt.Sprintf(
		"product:%d",
		productID,
	)
}
