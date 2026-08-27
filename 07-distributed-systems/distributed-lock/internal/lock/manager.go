package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	client *redis.Client
	config Config
}

type Lease struct {
	manager *Manager
	key     string
	token   string
}

func NewManager(
	client *redis.Client,
	config Config,
) (*Manager, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Manager{
		client: client,
		config: config,
	}, nil
}

func (m *Manager) Acquire(
	ctx context.Context,
	resource string,
) (*Lease, error) {
	if resource == "" {
		return nil, errors.New("lock resource is required")
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate lock token: %w", err)
	}

	key := "lock:" + resource

	for {
		acquired, err := m.client.SetNX(
			ctx,
			key,
			token,
			m.config.TTL,
		).Result()
		if err != nil {
			return nil, fmt.Errorf("acquire redis lock: %w", err)
		}

		if acquired {
			return &Lease{
				manager: m,
				key:     key,
				token:   token,
			}, nil
		}

		if err := wait(
			ctx,
			m.config.RetryInterval,
		); err != nil {
			return nil, errors.Join(
				ErrNotAcquired,
				err,
			)
		}
	}
}

func generateToken() (string, error) {
	value := make([]byte, 16)

	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(value), nil
}

func wait(
	ctx context.Context,
	duration time.Duration,
) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
