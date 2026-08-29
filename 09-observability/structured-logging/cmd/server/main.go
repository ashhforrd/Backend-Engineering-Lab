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

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/httpmiddleware"
	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/logging"
	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/order"
)

func main() {
	logger := logging.New(
		os.Stdout,
		logging.Config{
			Service:     "order-api",
			Environment: "local",
			Level:       slog.LevelInfo,
			AddSource:   false,
		},
	)

	slog.SetDefault(logger)

	orderService := order.NewService()
	apiHandler := api.NewHandler(orderService)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/orders",
		apiHandler.CreateOrder,
	)

	requestLogger :=
		httpmiddleware.NewRequestLogger()

	requestID :=
		httpmiddleware.NewRequestID(logger)

	handler := requestID.Wrap(
		requestLogger.Wrap(mux),
	)

	server := &http.Server{
		Addr:              ":8090",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"server started",
			"address",
			"http://localhost:8090",
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

	signalContext, stopSignals :=
		signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
	defer stopSignals()

	<-signalContext.Done()

	shutdownContext, cancelShutdown :=
		context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
	defer cancelShutdown()

	logger.Info("server shutdown started")

	if err := server.Shutdown(
		shutdownContext,
	); err != nil {
		logger.Error(
			"server shutdown failed",
			"error",
			err,
		)
	}

	logger.Info("server stopped")
}
