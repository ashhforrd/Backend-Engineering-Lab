package httpmiddleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/logging"
)

const RequestIDHeader = "X-Request-ID"

type RequestID struct {
	logger *slog.Logger
}

func NewRequestID(
	logger *slog.Logger,
) *RequestID {
	return &RequestID{
		logger: logger,
	}
}

func (m *RequestID) Wrap(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			requestID, err := generateRequestID()
			if err != nil {
				http.Error(
					writer,
					"generate request id",
					http.StatusInternalServerError,
				)
				return
			}

			requestLogger := m.logger.With(
				"request_id",
				requestID,
			)

			ctx := logging.WithLogger(
				request.Context(),
				requestLogger,
			)

			writer.Header().Set(
				RequestIDHeader,
				requestID,
			)

			next.ServeHTTP(
				writer,
				request.WithContext(ctx),
			)
		},
	)
}

func generateRequestID() (string, error) {
	value := make([]byte, 16)

	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(value), nil
}
