package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestLoggerRedactsSensitiveAttributes(t *testing.T) {
	var output bytes.Buffer

	logger := New(
		&output,
		Config{
			Service:     "test-service",
			Environment: "test",
			Level:       slog.LevelInfo,
		},
	)

	logger.Info(
		"test message",
		"email",
		"private@example.com",
		"payment_token",
		"secret-token",
		slog.Group(
			"credentials",
			"password",
			"secret-password",
		),
	)

	logLine := output.String()

	for _, secret := range []string{
		"private@example.com",
		"secret-token",
		"secret-password",
	} {
		if strings.Contains(logLine, secret) {
			t.Fatalf(
				"log contains sensitive value %q: %s",
				secret,
				logLine,
			)
		}
	}

	if count := strings.Count(
		logLine,
		redactedValue,
	); count != 3 {
		t.Fatalf(
			"expected 3 redacted values, got %d: %s",
			count,
			logLine,
		)
	}
}
