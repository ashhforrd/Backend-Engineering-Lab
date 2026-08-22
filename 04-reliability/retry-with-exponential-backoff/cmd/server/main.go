package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/payment"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/retry"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/simulator"
)

func main() {
	const address = ":8083"

	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	downstreamClient := downstream.NewClient(
		"http://localhost"+address+"/downstream",
		httpClient,
	)

	retryPolicy := retry.Policy{
		MaxAttempts:  4,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2,
		Jitter:       0.2,
	}

	paymentService := payment.NewService(
		downstreamClient,
		retryPolicy,
	)

	apiHandler := api.NewHandler(paymentService)
	simulatorHandler := simulator.NewHandler()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/payments",
		apiHandler.CreatePayment,
	)

	mux.HandleFunc(
		"POST /downstream/payments",
		simulatorHandler.CreatePayment,
	)

	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on http://localhost%s", address)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start server: %v", err)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-shutdownSignal

	log.Println("shutting down server")

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
