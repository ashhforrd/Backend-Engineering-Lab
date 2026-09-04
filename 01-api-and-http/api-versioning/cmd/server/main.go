package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	v1 "github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi/v1"
	v2 "github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi/v2"
	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/user"
)

func main() {
	repository := user.NewMemoryRepository(
		[]user.User{
			{
				ID:        "user-001",
				FirstName: "Ada",
				LastName:  "Lovelace",
				Email:     "ada@example.com",
				Status:    user.StatusActive,
				CreatedAt: time.Date(
					2026,
					time.September,
					4,
					12,
					0,
					0,
					0,
					time.UTC,
				),
			},
		},
	)

	service := user.NewService(repository)
	mux := http.NewServeMux()

	v1.NewHandler(service).RegisterRoutes(mux)
	v2.NewHandler(service).RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf(
		"API versioning server listening on %s",
		server.Addr,
	)

	err := server.ListenAndServe()
	if err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
