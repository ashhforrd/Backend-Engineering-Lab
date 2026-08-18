package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestCreatePayment_ReplaysResponse(t *testing.T) {
	server := NewServer()
	handler := server.routes()

	body := `{"amount":1099,"currency":"USD"}`

	firstRequest := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader(body),
	)
	firstRequest.Header.Set("Idempotency-Key", "payment-001")

	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, firstRequest)

	secondRequest := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader(body),
	)
	secondRequest.Header.Set("Idempotency-Key", "payment-001")

	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)

	if firstResponse.Code != http.StatusCreated {
		t.Fatalf("expected first status 201, got %d", firstResponse.Code)
	}

	if secondResponse.Code != http.StatusCreated {
		t.Fatalf("expected second status 201, got %d", secondResponse.Code)
	}

	if firstResponse.Body.String() != secondResponse.Body.String() {
		t.Fatal("expected retry to return the original response")
	}

	if secondResponse.Header().Get("Idempotency-Replayed") != "true" {
		t.Fatal("expected retry to be marked as replayed")
	}
}

func TestCreatePayment_RejectsDifferentPayload(t *testing.T) {
	server := NewServer()
	handler := server.routes()

	firstRequest := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader(`{"amount":1099,"currency":"USD"}`),
	)
	firstRequest.Header.Set("Idempotency-Key", "payment-001")

	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, firstRequest)

	secondRequest := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader(`{"amount":2099,"currency":"USD"}`),
	)
	secondRequest.Header.Set("Idempotency-Key", "payment-001")

	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)

	if firstResponse.Code != http.StatusCreated {
		t.Fatalf("expected first status 201, got %d", firstResponse.Code)
	}

	if secondResponse.Code != http.StatusConflict {
		t.Fatalf("expected second status 409, got %d", secondResponse.Code)
	}
}

func TestCreatePayment_ConcurrentRequestsCreateOnePayment(t *testing.T) {
	server := NewServer()
	handler := server.routes()

	const requestCount = 20

	responses := make([]*httptest.ResponseRecorder, requestCount)
	start := make(chan struct{})

	var waitGroup sync.WaitGroup
	waitGroup.Add(requestCount)

	for i := 0; i < requestCount; i++ {
		go func(index int) {
			defer waitGroup.Done()

			request := httptest.NewRequest(
				http.MethodPost,
				"/payments",
				strings.NewReader(`{"amount":1099,"currency":"USD"}`),
			)
			request.Header.Set("Idempotency-Key", "payment-001")

			response := httptest.NewRecorder()
			responses[index] = response

			<-start
			handler.ServeHTTP(response, request)
		}(i)
	}

	close(start)
	waitGroup.Wait()

	expectedBody := responses[0].Body.String()

	for _, response := range responses {
		if response.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", response.Code)
		}

		if response.Body.String() != expectedBody {
			t.Fatal("expected all requests to return the same payment")
		}
	}

	if server.nextPaymentID != 1 {
		t.Fatalf("expected one payment, got %d", server.nextPaymentID)
	}
}
