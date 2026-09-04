package user

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	FindByID(ctx context.Context, id string) (User, error)
}

type MemoryRepository struct {
	users map[string]User
}

func NewMemoryRepository(
	users []User,
) *MemoryRepository {
	storedUsers := make(map[string]User, len(users))

	for _, currentUser := range users {
		storedUsers[currentUser.ID] = currentUser
	}

	return &MemoryRepository{
		users: storedUsers,
	}
}

func (r *MemoryRepository) FindByID(
	ctx context.Context,
	id string,
) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	foundUser, exists := r.users[id]
	if !exists {
		return User{}, ErrNotFound
	}

	return foundUser, nil
}
