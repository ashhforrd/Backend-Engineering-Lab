package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-rate-limit-headers/internal/ratelimit"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

func TestRateLimitMiddlewareCommunicatesQuota(t *testing.T) {
	limiter, err := ratelimit.NewLimiter(
		ratelimit.Config{Limit: 2, Window: time.Minute},
		fixedClock{now: time.Unix(1_000, 0)},
	)
	if err != nil {
		t.Fatalf("create limiter: %v", err)
	}

	handler := RateLimitMiddleware(limiter, MessageHandler())

	first := performRequest(handler)
	second := performRequest(handler)
	third := performRequest(handler)

	if first.Code != http.StatusOK || first.Header().Get("X-RateLimit-Remaining") != "1" {
		t.Fatalf("unexpected first response: status=%d headers=%v", first.Code, first.Header())
	}

	if second.Code != http.StatusOK || second.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Fatalf("unexpected second response: status=%d headers=%v", second.Code, second.Header())
	}

	if third.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", third.Code)
	}

	if third.Header().Get("Retry-After") != "60" {
		t.Fatalf("expected Retry-After 60, got %q", third.Header().Get("Retry-After"))
	}

	if third.Header().Get("X-RateLimit-Limit") != "2" {
		t.Fatalf("expected limit header 2, got %q", third.Header().Get("X-RateLimit-Limit"))
	}
}

func performRequest(handler http.Handler) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	request.Header.Set("X-API-Key", "client-001")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	return response
}
