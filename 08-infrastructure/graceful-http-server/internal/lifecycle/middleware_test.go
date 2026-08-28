package lifecycle

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareDrainsActiveRequestAndRejectsNewOnes(
	t *testing.T,
) {
	state := NewState()
	middleware := NewMiddleware(state)

	started := make(chan struct{})
	release := make(chan struct{})

	handler := middleware.Track(
		http.HandlerFunc(
			func(
				writer http.ResponseWriter,
				request *http.Request,
			) {
				close(started)
				<-release
				writer.WriteHeader(http.StatusOK)
			},
		),
	)

	firstResponse := httptest.NewRecorder()
	firstDone := make(chan struct{})

	go func() {
		defer close(firstDone)

		handler.ServeHTTP(
			firstResponse,
			httptest.NewRequest(
				http.MethodGet,
				"/api/work",
				nil,
			),
		)
	}()

	<-started

	if active := middleware.ActiveRequests(); active != 1 {
		t.Fatalf(
			"expected 1 active request, got %d",
			active,
		)
	}

	state.BeginShutdown()

	secondResponse := httptest.NewRecorder()

	handler.ServeHTTP(
		secondResponse,
		httptest.NewRequest(
			http.MethodGet,
			"/api/work",
			nil,
		),
	)

	if secondResponse.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			secondResponse.Code,
		)
	}

	close(release)
	<-firstDone

	if firstResponse.Code != http.StatusOK {
		t.Fatalf(
			"expected active request status %d, got %d",
			http.StatusOK,
			firstResponse.Code,
		)
	}

	if active := middleware.ActiveRequests(); active != 0 {
		t.Fatalf(
			"expected 0 active requests, got %d",
			active,
		)
	}
}
