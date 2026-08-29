package httpmiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/logging"
)

func TestRequestIDCorrelatesBusinessAndHTTPLogs(
	t *testing.T,
) {
	var output bytes.Buffer

	logger := logging.New(
		&output,
		logging.Config{
			Service:     "test-service",
			Environment: "test",
			Level:       slog.LevelInfo,
		},
	)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /test",
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			logging.FromContext(
				request.Context(),
			).InfoContext(
				request.Context(),
				"business event",
			)

			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte("created"))
		},
	)

	handler := NewRequestID(logger).Wrap(
		NewRequestLogger().Wrap(mux),
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	requestID := response.Header().Get(
		RequestIDHeader,
	)
	if requestID == "" {
		t.Fatal("expected response request ID")
	}

	decoder := json.NewDecoder(&output)
	logCount := 0

	for decoder.More() {
		var record map[string]any

		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode log: %v", err)
		}

		if record["request_id"] != requestID {
			t.Fatalf(
				"expected request ID %q, got %v",
				requestID,
				record["request_id"],
			)
		}

		logCount++
	}

	if logCount != 2 {
		t.Fatalf("expected 2 log records, got %d", logCount)
	}
}

func TestResponseWriterKeepsFirstStatusCode(
	t *testing.T,
) {
	recorder := httptest.NewRecorder()
	writer := NewResponseWriter(recorder)

	writer.WriteHeader(http.StatusCreated)
	writer.WriteHeader(http.StatusInternalServerError)

	if writer.StatusCode() != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			writer.StatusCode(),
		)
	}
}

func TestContextFallbackLoggerIsAvailable(
	t *testing.T,
) {
	if logging.FromContext(
		context.Background(),
	) == nil {
		t.Fatal("expected fallback logger")
	}
}
