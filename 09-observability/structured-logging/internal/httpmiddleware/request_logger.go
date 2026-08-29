package httpmiddleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/09-observability/structured-logging/internal/logging"
)

type RequestLogger struct{}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

func (m *RequestLogger) Wrap(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			startedAt := time.Now()
			responseWriter := NewResponseWriter(writer)

			next.ServeHTTP(
				responseWriter,
				request,
			)

			duration := time.Since(startedAt)
			logger := logging.FromContext(
				request.Context(),
			)

			level := levelForStatus(
				responseWriter.StatusCode(),
			)

			logger.LogAttrs(
				request.Context(),
				level,
				"http request completed",
				slog.String(
					"http.method",
					request.Method,
				),
				slog.String(
					"http.route",
					request.Pattern,
				),
				slog.String(
					"http.path",
					request.URL.Path,
				),
				slog.Int(
					"http.status_code",
					responseWriter.StatusCode(),
				),
				slog.Int64(
					"http.duration_ms",
					duration.Milliseconds(),
				),
				slog.Int(
					"http.response_bytes",
					responseWriter.BytesWritten(),
				),
				slog.String(
					"client.address",
					clientAddress(request),
				),
				slog.String(
					"user_agent",
					request.UserAgent(),
				),
			)
		},
	)
}

func levelForStatus(
	statusCode int,
) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError

	case statusCode >= 400:
		return slog.LevelWarn

	default:
		return slog.LevelInfo
	}
}

func clientAddress(
	request *http.Request,
) string {
	host, _, err := net.SplitHostPort(
		request.RemoteAddr,
	)
	if err != nil {
		return request.RemoteAddr
	}

	return host
}
