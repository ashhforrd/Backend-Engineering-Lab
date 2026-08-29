# Structured Logging

## Problem

Plain-text logs are difficult to search, aggregate, and correlate across request
flows. Logs also risk exposing sensitive data when fields are added without a
consistent policy.

## Design

```text
HTTP Request
  ↓
Request ID Middleware
  ↓
Request Logging Middleware
  ↓
Order Handler
  ↓
Order Service
```

Every log in one request carries the same `request_id`. JSON fields allow log
platforms to filter by level, service, status, error type, and duration.

## Implementation

The project uses Go's `log/slog`.

- `internal/logging` configures JSON output, context propagation, and redaction.
- `internal/httpmiddleware` assigns request IDs and records HTTP outcomes.
- `internal/order` emits business-event logs.
- `internal/api` translates HTTP requests into order operations.
- `cmd/server` wires middleware and handles shutdown.

Sensitive keys such as email, payment tokens, authorization, and passwords are
automatically replaced with `[REDACTED]`.

## Failure Cases

- Invalid requests generate `WARN` HTTP logs.
- Server errors generate `ERROR` HTTP logs.
- Business errors include a stable `error_type`.
- Missing context falls back to the default logger.
- Sensitive values must never be used as log messages.
- Request bodies and authorization headers are not logged.
- High-cardinality fields can increase logging cost.

## What I Learned

- Structured logs store fields instead of embedding data in messages.
- Request IDs correlate HTTP and business logs.
- Context carries request-scoped logging metadata.
- Log levels describe operational severity.
- Redaction provides defense in depth.
- Completion logs capture status, duration, and response size.
- Consistent field names make logs machine-searchable.

## Running

Start the server:

```bash
go run ./cmd/server
```

Create an order:

```bash
curl -X POST http://localhost:8090/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "customer-123",
    "email": "student@example.com",
    "amount": 250000,
    "paymentToken": "secret-token",
    "simulateFailure": false
  }'
```

Filter error logs with `jq`:

```bash
go run ./cmd/server |
  jq 'select(.level == "ERROR")'
```

Run tests:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```