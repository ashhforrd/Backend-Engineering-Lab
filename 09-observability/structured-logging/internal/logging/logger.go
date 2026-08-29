package logging

import (
	"io"
	"log/slog"
)

type Config struct {
	Service     string
	Environment string
	Level       slog.Level
	AddSource   bool
}

func New(
	output io.Writer,
	config Config,
) *slog.Logger {
	handler := slog.NewJSONHandler(
		output,
		&slog.HandlerOptions{
			Level:     config.Level,
			AddSource: config.AddSource,
			ReplaceAttr: func(
				groups []string,
				attribute slog.Attr,
			) slog.Attr {
				switch attribute.Key {
				case slog.TimeKey:
					attribute.Key = "timestamp"

				case slog.MessageKey:
					attribute.Key = "message"
				}

				return attribute

			},
		},
	)

	redactingHandler := NewRedactingHandler(
		handler,
		[]string{
			"email",
			"payment_token",
			"authorization",
			"password",
		},
	)

	return slog.New(redactingHandler).With(
		"service",
		config.Service,
		"environment",
		config.Environment,
	)
}
