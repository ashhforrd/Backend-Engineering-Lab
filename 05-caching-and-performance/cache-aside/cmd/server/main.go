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

	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/postgres"
	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/product"
	redisstore "github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	databaseURL := envOrDefault(
		"DATABASE_URL",
		"postgres://postgres:postgres@localhost:5435/cache_aside",
	)

	redisAddress := envOrDefault(
		"REDIS_ADDRESS",
		"localhost:6380",
	)

	startupContext, cancelStartup := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancelStartup()

	databasePool, err := pgxpool.New(
		startupContext,
		databaseURL,
	)
	if err != nil {
		logger.Error(
			"create database pool",
			"error",
			err,
		)
		os.Exit(1)
	}
	defer databasePool.Close()

	if err := databasePool.Ping(startupContext); err != nil {
		logger.Error(
			"connect to database",
			"error",
			err,
		)
		os.Exit(1)
	}

	redisClient := goredis.NewClient(
		&goredis.Options{
			Addr: redisAddress,
		},
	)
	defer redisClient.Close()

	repository := postgres.NewRepository(databasePool)
	cache := redisstore.NewCache(redisClient)

	productService := product.NewService(
		repository,
		cache,
		30*time.Second,
		logger,
	)

	apiHandler := api.NewHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/products/{id}",
		apiHandler.GetProduct,
	)

	mux.HandleFunc(
		"PUT /api/products/{id}",
		apiHandler.UpdateProduct,
	)

	server := &http.Server{
		Addr:              ":8085",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"address",
			"http://localhost:8085",
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
			"graceful shutdown failed",
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
