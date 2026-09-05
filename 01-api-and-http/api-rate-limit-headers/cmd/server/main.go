package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-rate-limit-headers/internal/httpapi"
	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-rate-limit-headers/internal/ratelimit"
)

func main() {
	limiter, err := ratelimit.NewLimiter(
		ratelimit.Config{
			Limit:  5,
			Window: 30 * time.Second,
		},
		ratelimit.SystemClock{},
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle(
		"GET /api/message",
		httpapi.RateLimitMiddleware(
			limiter,
			httpapi.MessageHandler(),
		),
	)

	mux.HandleFunc(
		"GET /health",
		func(
			writer http.ResponseWriter,
			_ *http.Request,
		) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("ok\n"))
		},
	)

	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
