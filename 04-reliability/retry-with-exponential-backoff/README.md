# Retry with Exponential Backoff

## Problem

Downstream services can fail temporarily because of overload, network errors,
timeouts, or deployments. Immediately failing every request reduces reliability,
while retrying continuously can make the outage worse.

## Design

The application calls a simulated payment provider through a reusable retry
executor.

```text
Client
  -> Payment API
  -> Retry Executor
  -> Downstream HTTP Client
  -> Payment Simulator
```

Each failed retry waits longer than the previous one:

```text
delay = initialDelay * multiplier^(attempt - 1)
```

Jitter adds randomness so multiple clients do not retry simultaneously.

## Implementation

The project uses Go's standard library.

- `internal/retry` implements policy validation, backoff, and retry execution.
- `internal/downstream` calls the downstream HTTP service.
- `internal/payment` coordinates payment processing and retry.
- `internal/simulator` simulates temporary downstream failures.
- `internal/api` exposes the public HTTP endpoint.
- `cmd/server` wires dependencies and runs the server.

Only temporary failures such as `408`, `429`, and selected `5xx` responses are
retried. Permanent client errors are returned immediately.

## Failure Cases

- A cancelled context stops the retry loop.
- A permanent error is not retried.
- Retry stops after `MaxAttempts`.
- Backoff never exceeds `MaxDelay`.
- HTTP calls have a per-attempt timeout.
- Jitter reduces synchronized retries during outages.

Retrying a non-idempotent operation can create duplicate effects. A production
payment API should also use an idempotency key.

## What I Learned

- Functions can be passed as values in Go.
- Generics make the retry executor reusable for different result types.
- Closures can capture request state and attempt counters.
- Context supports cancellation and deadlines.
- Exponential backoff prevents aggressive retry loops.
- Jitter helps avoid the thundering herd problem.
- Retry decisions must distinguish temporary and permanent failures.

## Running

Run the tests:

```bash
go test ./...
```

Start the application:

```bash
go run ./cmd/server
```

Test a request that fails twice before succeeding:

```bash
curl -i -X POST http://localhost:8083/api/payments \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "order-001",
    "amount": 150000,
    "failuresBeforeSuccess": 2
  }'
```

Expected response:

```json
{
  "paymentId": "pay-order-001",
  "status": "SUCCESS",
  "attempts": 3
}
```

Use a new `orderId` for each experiment because the simulator stores attempt
counts in memory.