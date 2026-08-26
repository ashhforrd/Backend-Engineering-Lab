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

	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/job"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/producer"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/queue"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/tasks"
	"github.com/ashhforrd/backend-engineering-lab/06-async-processing/basic-job-queue/internal/worker"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	jobStore := job.NewStore()
	jobQueue := queue.New(100)

	taskHandlers := tasks.NewHandlers(logger)

	dispatcher := worker.NewDispatcher(
		map[job.Type]worker.Handler{
			job.TypeSendEmail:      taskHandlers.SendEmail,
			job.TypeGenerateReport: taskHandlers.GenerateReport,
		},
	)

	workerPool := worker.NewPool(
		3,
		jobQueue,
		jobStore,
		dispatcher,
		logger,
	)

	workerContext, cancelWorkers := context.WithCancel(
		context.Background(),
	)
	defer cancelWorkers()

	workerPool.Start(workerContext)

	producerService := producer.NewService(
		jobStore,
		jobQueue,
	)

	apiHandler := api.NewHandler(
		producerService,
		jobStore,
		jobQueue,
	)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/jobs",
		apiHandler.CreateJob,
	)

	mux.HandleFunc(
		"GET /api/jobs/{id}",
		apiHandler.GetJob,
	)

	mux.HandleFunc(
		"GET /api/queue",
		apiHandler.GetQueueStatus,
	)

	server := &http.Server{
		Addr:              ":8086",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"address",
			"http://localhost:8086",
			"workers",
			3,
			"queue_capacity",
			jobQueue.Capacity(),
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

	cancelWorkers()
	workerPool.Wait()

	logger.Info("application stopped")
}
