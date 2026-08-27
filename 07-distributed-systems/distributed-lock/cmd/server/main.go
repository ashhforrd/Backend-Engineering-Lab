package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/critical"
	"github.com/ashhforrd/backend-engineering-lab/07-distributed-systems/distributed-lock/internal/lock"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	instanceID := envOrDefault(
		"INSTANCE_ID",
		"instance-1",
	)

	serverAddress := envOrDefault(
		"SERVER_ADDRESS",
		":8088",
	)

	redisAddress := envOrDefault(
		"REDIS_ADDRESS",
		"localhost:6381",
	)

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: redisAddress,
		},
	)
	defer redisClient.Close()

	startupContext, cancelStartup := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancelStartup()

	if err := redisClient.Ping(startupContext).Err(); err != nil {
		logger.Error(
			"connect to redis",
			"error",
			err,
		)
		os.Exit(1)
	}

	lockManager, err := lock.NewManager(
		redisClient,
		lock.Config{
			TTL:           3 * time.Second,
			RetryInterval: 100 * time.Millisecond,
		},
	)
	if err != nil {
		logger.Error(
			"create lock manager",
			"error",
			err,
		)
		os.Exit(1)
	}

	criticalService := critical.NewService(
		lockManager,
		redisClient,
		instanceID,
		2*time.Second,
		logger,
	)

	apiHandler := api.NewHandler(criticalService)
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/critical/{resource}",
		apiHandler.Execute,
	)

	server := &http.Server{
		Addr:              serverAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"instance_id",
			instanceID,
			"address",
			serverAddress,
		)

		if err := server.ListenAndServe(); !errors.Is(
			err,
			http.ErrServerClosed,
		) {
			logger.Error(
				"server failed",
				"error",
				err,
			)
			os.Exit(1)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-shutdownSignal

	shutdownContext, cancelShutdown := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error(
			"server shutdown failed",
			"error",
			err,
		)
	}
}

func envOrDefault(
	key string,
	fallback string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
