package logging

import (
	"context"
	"log/slog"
	"strings"
)

const redactedValue = "[REDACTED]"

type RedactingHandler struct {
	next          slog.Handler
	sensitiveKeys map[string]struct{}
}

func NewRedactingHandler(
	next slog.Handler,
	keys []string,
) *RedactingHandler {
	sensitiveKeys := make(
		map[string]struct{},
		len(keys),
	)

	for _, key := range keys {
		sensitiveKeys[strings.ToLower(key)] =
			struct{}{}
	}

	return &RedactingHandler{
		next:          next,
		sensitiveKeys: sensitiveKeys,
	}
}

func (h *RedactingHandler) Enabled(
	ctx context.Context,
	level slog.Level,
) bool {
	return h.next.Enabled(ctx, level)
}

func (h *RedactingHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {
	redactedRecord := slog.NewRecord(
		record.Time,
		record.Level,
		record.Message,
		record.PC,
	)

	record.Attrs(
		func(attribute slog.Attr) bool {
			redactedRecord.AddAttrs(
				h.redact(attribute),
			)

			return true
		},
	)

	return h.next.Handle(
		ctx,
		redactedRecord,
	)
}

func (h *RedactingHandler) WithAttrs(
	attributes []slog.Attr,
) slog.Handler {
	redacted := make(
		[]slog.Attr,
		0,
		len(attributes),
	)

	for _, attribute := range attributes {
		redacted = append(
			redacted,
			h.redact(attribute),
		)
	}

	return &RedactingHandler{
		next:          h.next.WithAttrs(redacted),
		sensitiveKeys: h.sensitiveKeys,
	}
}

func (h *RedactingHandler) WithGroup(
	name string,
) slog.Handler {
	return &RedactingHandler{
		next:          h.next.WithGroup(name),
		sensitiveKeys: h.sensitiveKeys,
	}
}

func (h *RedactingHandler) redact(
	attribute slog.Attr,
) slog.Attr {
	attribute.Value = attribute.Value.Resolve()

	if _, sensitive := h.sensitiveKeys[strings.ToLower(attribute.Key)]; sensitive {
		return slog.String(
			attribute.Key,
			redactedValue,
		)
	}

	if attribute.Value.Kind() == slog.KindGroup {
		group := attribute.Value.Group()

		for index := range group {
			group[index] = h.redact(group[index])
		}

		attribute.Value = slog.GroupValue(group...)
	}

	return attribute
}
