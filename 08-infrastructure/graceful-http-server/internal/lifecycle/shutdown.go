package lifecycle

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ShutdownCoordinator struct {
	server           *http.Server
	state            *State
	cancelBackground context.CancelFunc
	drainDelay       time.Duration
	shutdownTimeout  time.Duration
	logger           *slog.Logger
}

func NewShutdownCoordinator(
	server *http.Server,
	state *State,
	cancelBackground context.CancelFunc,
	drainDelay time.Duration,
	shutdownTimeout time.Duration,
	logger *slog.Logger,
) *ShutdownCoordinator {
	return &ShutdownCoordinator{
		server:           server,
		state:            state,
		cancelBackground: cancelBackground,
		drainDelay:       drainDelay,
		shutdownTimeout:  shutdownTimeout,
		logger:           logger,
	}
}

func (c *ShutdownCoordinator) Shutdown() error {
	c.logger.Info(
		"shutdown started",
	)

	c.state.BeginShutdown()

	c.logger.Info(
		"server marked not ready",
		"drain_delay",
		c.drainDelay,
	)

	time.Sleep(c.drainDelay)

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		c.shutdownTimeout,
	)
	defer cancel()

	err := c.server.Shutdown(shutdownContext)
	if err != nil {
		c.logger.Warn(
			"graceful shutdown deadline exceeded",
			"error",
			err,
		)

		if closeErr := c.server.Close(); closeErr != nil {
			return fmt.Errorf(
				"force close server: %w",
				closeErr,
			)
		}
	}

	c.cancelBackground()

	c.logger.Info(
		"shutdown completed",
	)

	return err
}
