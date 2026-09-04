package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi/v1"
	v2 "github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi/v2"
	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/user"
)

func newTestHandler() http.Handler {
	repository := user.NewMemoryRepository([]user.User{
		{
			ID:        "user-001",
			FirstName: "Ada",
			LastName:  "Lovelace",
			Email:     "ada@example.com",
			Status:    user.StatusActive,
			CreatedAt: time.Date(2026, time.September, 4, 12, 0, 0, 0, time.UTC),
		},
	})

	service := user.NewService(repository)
	mux := http.NewServeMux()

	v1.NewHandler(service).RegisterRoutes(mux)
	v2.NewHandler(service).RegisterRoutes(mux)

	return mux
}

func TestV1PreservesTheLegacyContract(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-001", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["name"] != "Ada Lovelace" {
		t.Fatalf("expected legacy name field, got %#v", body["name"])
	}

	if _, exists := body["profile"]; exists {
		t.Fatal("v1 must not expose the v2 profile field")
	}

	if _, exists := body["status"]; exists {
		t.Fatal("v1 must not expose the v2 status field")
	}

	if response.Header().Get("Deprecation") != "@1788480000" {
		t.Fatal("v1 must include the deprecation header")
	}

	if response.Header().Get("Sunset") == "" {
		t.Fatal("v1 must include the sunset header")
	}

	if response.Header().Get("Link") != "</api/v2/users/user-001>; rel=\"successor-version\"" {
		t.Fatal("v1 must link to the successor version")
	}
}

func TestV2ReturnsTheEvolvedContract(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v2/users/user-001", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	profile, ok := body["profile"].(map[string]any)
	if !ok {
		t.Fatal("v2 must expose a profile object")
	}

	if profile["firstName"] != "Ada" || profile["lastName"] != "Lovelace" {
		t.Fatalf("unexpected v2 profile: %#v", profile)
	}

	if body["status"] != "ACTIVE" {
		t.Fatalf("expected ACTIVE status, got %#v", body["status"])
	}

	if _, exists := body["name"]; exists {
		t.Fatal("v2 must not expose the legacy name field")
	}

	if response.Header().Get("Deprecation") != "" {
		t.Fatal("v2 must not be marked as deprecated")
	}
}

func TestBothVersionsReturnNotFound(t *testing.T) {
	paths := []string{
		"/api/v1/users/missing",
		"/api/v2/users/missing",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()

			newTestHandler().ServeHTTP(response, request)

			if response.Code != http.StatusNotFound {
				t.Fatalf("expected status 404, got %d", response.Code)
			}
		})
	}
}
