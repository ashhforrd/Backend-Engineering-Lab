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

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/api"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/circuitbreaker"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/product"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/circuit-breaker/internal/simulator"
)

func main() {
	const address = ":8084"

	breaker, err := circuitbreaker.New(
		circuitbreaker.Config{
			FailureThreshold:     3,
			OpenTimeout:          5 * time.Second,
			HalfOpenMaxRequests:  1,
			HalfOpenSuccessLimit: 1,
		},
	)
	if err != nil {
		log.Fatalf("create circuit breaker: %v", err)
	}

	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	downstreamClient := downstream.NewClient(
		"http://localhost"+address+"/downstream",
		httpClient,
	)

	productService := product.NewService(
		downstreamClient,
		breaker,
	)

	apiHandler := api.NewHandler(
		productService,
		breaker,
	)

	simulatorHandler := simulator.NewHandler()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/products/{id}",
		apiHandler.GetProduct,
	)

	mux.HandleFunc(
		"GET /api/circuit-breaker",
		apiHandler.GetCircuitStatus,
	)

	mux.HandleFunc(
		"GET /downstream/products/{id}",
		simulatorHandler.GetProduct,
	)

	mux.HandleFunc(
		"PUT /simulator/availability",
		simulatorHandler.SetAvailability,
	)

	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf(
			"server listening on http://localhost%s",
			address,
		)

		if err := server.ListenAndServe(); !errors.Is(
			err,
			http.ErrServerClosed,
		) {
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

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}