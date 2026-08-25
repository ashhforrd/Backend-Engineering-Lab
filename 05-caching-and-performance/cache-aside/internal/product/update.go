package product

import (
	"context"
)

func (s *Service) Update(
	ctx context.Context,
	product Product,
) (Product, error) {
	updatedProduct, err := s.repository.Update(
		ctx,
		product,
	)
	if err != nil {
		return Product{}, err
	}

	if err := s.cache.Delete(
		ctx,
		updatedProduct.ID,
	); err != nil {
		s.logger.WarnContext(
			ctx,
			"cache invalid failed",
			"product_id",
			updatedProduct.ID,
			"error",
			err,
		)
	}

	return updatedProduct, nil
}
