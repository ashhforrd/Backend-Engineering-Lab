package product

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type Service struct {
	repository Repository
	cache      Cache
	cacheTTL   time.Duration
	logger     *slog.Logger
}

func NewService(
	repository Repository,
	cache Cache,
	cacheTTL time.Duration,
	logger *slog.Logger,
) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
		cacheTTL:   cacheTTL,
		logger:     logger,
	}
}

func (s *Service) Get(
	ctx context.Context,
	productID int64,
) (Result, error) {
	cacheProduct, err := s.cache.Get(
		ctx,
		productID,
	)
	if err == nil {
		return Result{
			Product: cacheProduct,
			Source:  SourceCache,
		}, nil
	}

	if !errors.Is(err, ErrCacheMiss) {
		s.logger.WarnContext(
			ctx,
			"cache read failed",
			"product_id",
			productID,
			"error",
			err,
		)
	}

	databaseProduct, err := s.repository.FindByID(
		ctx,
		productID,
	)
	if err != nil {
		return Result{}, err
	}

	if err := s.cache.Set(
		ctx,
		databaseProduct,
		s.cacheTTL,
	); err != nil {
		s.logger.WarnContext(
			ctx,
			"cache write failed",
			"product_id",
			productID,
			"error",
			err,
		)
	}

	return Result{
		Product: databaseProduct,
		Source:  SourceDatabase,
	}, nil
}
