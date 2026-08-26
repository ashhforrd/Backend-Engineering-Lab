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

	"github.com/ashhforrd/backend-engineering-lab/03-concurrency/bounded-worker-pool/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/03-concurrency/bounded-worker-pool/internal/pool"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	workerPool, err := pool.New(
		pool.Config{
			WorkerCount:   3,
			QueueCapacity: 10,
		},
	)
	if err != nil {
		logger.Info(
			"create worker pool",
			"error",
			err,
		)
		os.Exit(1)
	}

	if err := workerPool.Start(
		context.Background(),
	); err != nil {
		logger.Info(
			"start worker pool",
			"error",
			"err",
		)
		os.Exit(1)
	}

	go func() {
		for result := range workerPool.Results() {
			logger.Info(
				"task finished",
				"task_id",
				result.TaskID,
				"worker_id",
				result.WorkerID,
				"duration",
				result.Duration(),
				"error",
				result.Err,
			)
		}
	}()

	apiHandler := api.NewHandler(workerPool)
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/tasks",
		apiHandler.SubmitTasks,
	)

	mux.HandleFunc(
		"GET /api/stats",
		apiHandler.GetStats,
	)

	server := &http.Server{
		Addr:              ":8087",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"address",
			"http://localhost:8087",
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

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error(
			"shutdown server",
			"error",
			err,
		)
	}

	workerPool.Stop()

	logger.Info("application stopped")
}
