package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/downstream"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/payment"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/retry"
	"github.com/ashhforrd/backend-engineering-lab/04-reliability/retry-with-exponential-backoff/internal/simulator"
)

func newTestHandler(t *testing.T) (*Handler, func()) {
	t.Helper()

	simulatorHandler := simulator.NewHandler()

	downstreamServer := httptest.NewServer(
		http.HandlerFunc(simulatorHandler.CreatePayment),
	)

	client := downstream.NewClient(
		downstreamServer.URL,
		downstreamServer.Client(),
	)

	policy := retry.Policy{
		MaxAttempts:  4,
		InitialDelay: time.Millisecond,
		MaxDelay:     4 * time.Millisecond,
		Multiplier:   2,
		Jitter:       0,
	}

	service := payment.NewService(client, policy)

	return NewHandler(service), downstreamServer.Close
}

func TestCreatePaymentRetriesTemporaryFailures(t *testing.T) {
	handler, closeServer := newTestHandler(t)
	defer closeServer()

	requestBody := []byte(`{
		"orderId": "order-001",
		"amount": 150000,
		"failuresBeforeSuccess": 2
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/payments",
		bytes.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.CreatePayment(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var result payment.Result

	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", result.Attempts)
	}

	if result.Status != "SUCCESS" {
		t.Fatalf("expected SUCCESS, got %q", result.Status)
	}
}

func TestCreatePaymentFailsWhenAttemptsAreExhausted(t *testing.T) {
	handler, closeServer := newTestHandler(t)
	defer closeServer()

	requestBody := []byte(`{
		"orderId": "order-002",
		"amount": 150000,
		"failuresBeforeSuccess": 10
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/payments",
		bytes.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()

	handler.CreatePayment(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadGateway,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
