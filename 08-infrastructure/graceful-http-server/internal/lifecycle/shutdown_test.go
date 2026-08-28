package lifecycle

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestShutdownWaitsForActiveHandler(t *testing.T) {
	state := NewState()
	handlerStarted := make(chan struct{})
	releaseHandler := make(chan struct{})

	server := &http.Server{
		Handler: http.HandlerFunc(
			func(
				writer http.ResponseWriter,
				request *http.Request,
			) {
				close(handlerStarted)
				<-releaseHandler
				writer.WriteHeader(http.StatusOK)
			},
		),
	}

	listener, err := net.Listen(
		"tcp",
		"127.0.0.1:0",
	)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	serverStopped := make(chan error, 1)

	go func() {
		serverStopped <- server.Serve(listener)
	}()

	requestDone := make(chan error, 1)

	go func() {
		response, err := http.Get(
			"http://" + listener.Addr().String(),
		)
		if err == nil {
			_ = response.Body.Close()
		}

		requestDone <- err
	}()

	<-handlerStarted

	backgroundContext, cancelBackground :=
		context.WithCancel(context.Background())

	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	coordinator := NewShutdownCoordinator(
		server,
		state,
		cancelBackground,
		0,
		time.Second,
		logger,
	)

	shutdownDone := make(chan error, 1)

	go func() {
		shutdownDone <- coordinator.Shutdown()
	}()

	select {
	case err := <-shutdownDone:
		t.Fatalf(
			"shutdown returned before handler finished: %v",
			err,
		)

	case <-time.After(20 * time.Millisecond):
	}

	close(releaseHandler)

	if err := <-shutdownDone; err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	if err := <-requestDone; err != nil {
		t.Fatalf("active request failed: %v", err)
	}

	select {
	case <-backgroundContext.Done():

	default:
		t.Fatal("background context was not cancelled")
	}

	if err := <-serverStopped; !errors.Is(
		err,
		http.ErrServerClosed,
	) {
		t.Fatalf("unexpected server result: %v", err)
	}

	if state.IsReady() {
		t.Fatal("server must be not ready after shutdown")
	}
}
