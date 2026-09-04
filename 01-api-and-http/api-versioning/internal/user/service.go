package user

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidID = errors.New("user ID is required")

type Service struct {
	repository Repository
}

func NewService(
	repository Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (User, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return User{}, ErrInvalidID
	}

	return s.repository.FindByID(ctx, id)
}
