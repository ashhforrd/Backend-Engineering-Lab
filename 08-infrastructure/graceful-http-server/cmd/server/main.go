package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/08-infrastructure/graceful-http-server/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/08-infrastructure/graceful-http-server/internal/background"
	"github.com/ashhforrd/backend-engineering-lab/08-infrastructure/graceful-http-server/internal/lifecycle"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	state := lifecycle.NewState()
	middleware := lifecycle.NewMiddleware(state)
	apiHandler := api.NewHandler(state, middleware)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health/live",
		apiHandler.Liveness,
	)

	mux.HandleFunc(
		"GET /health/ready",
		apiHandler.Readiness,
	)

	mux.Handle(
		"GET /api/work",
		middleware.Track(
			http.HandlerFunc(
				apiHandler.SlowWork,
			),
		),
	)

	server := &http.Server{
		Addr:              ":8088",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	backgroundContext, cancelBackground :=
		context.WithCancel(context.Background())

	backgroundWorker := background.NewWorker(
		logger,
		2*time.Second,
	)

	var backgroundWaitGroup sync.WaitGroup
	backgroundWaitGroup.Add(1)

	go func() {
		defer backgroundWaitGroup.Done()
		backgroundWorker.Run(backgroundContext)
	}()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(
			"server started",
			"address",
			"http://localhost:8088",
		)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	signalContext, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	select {
	case err := <-serverErrors:
		logger.Error(
			"server failed",
			"error",
			err,
		)
		cancelBackground()

	case <-signalContext.Done():
		coordinator :=
			lifecycle.NewShutdownCoordinator(
				server,
				state,
				cancelBackground,
				2*time.Second,
				5*time.Second,
				logger,
			)

		if err := coordinator.Shutdown(); err != nil {
			logger.Error(
				"shutdown failed",
				"error",
				err,
			)
		}
	}

	backgroundWaitGroup.Wait()
}
