# Idempotency Key API

## Problem

A client may retry a request after a timeout because it does not know whether the server completed the operation.

Without idempotency handling, retrying a payment request could create multiple payments.

## Design

The client sends an `Idempotency-Key` header with each payment request.

The server stores:

- The idempotency key
- A hash of the request payload
- The HTTP status code
- The response body

Behavior:

- A new key creates a payment.
- The same key and payload return the stored response.
- The same key with a different payload returns `409 Conflict`.
- Concurrent requests with the same key create only one payment.

## Implementation

The project uses Go's standard HTTP library and an in-memory map.

A SHA-256 hash connects an idempotency key to its original payload. A mutex makes checking and storing a key atomic.

The implementation is separated into:

- `main.go`: application entry point
- `handler.go`: HTTP routing and request handling
- `payment.go`: payment domain types
- `idempotency.go`: idempotency record and request hashing
- `server.go`: shared server state
- `server_test.go`: behavior and concurrency tests

## Failure Cases

- Missing `Idempotency-Key` returns `400 Bad Request`.
- Invalid request JSON returns `400 Bad Request`.
- Invalid payment data returns `400 Bad Request`.
- Reusing a key with a different payload returns `409 Conflict`.
- Restarting the server removes all idempotency records.
- Multiple server instances do not share records.
- Records never expire, so memory usage can grow indefinitely.
- A global mutex serializes payment creation.

## What I Learned

- Client retries can duplicate non-idempotent operations.
- An idempotency key must be associated with the original payload.
- Checking and storing the key must be atomic.
- Returning the original status and body gives retries a consistent result.
- In-memory idempotency works for learning but not for distributed production systems.

## Running

Start the server:

```bash
go run .
```

Create a payment:

```bash
curl -i \
  -X POST http://localhost:8000/payments \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: payment-001" \
  -d '{"amount":1099,"currency":"USD"}'
```

Run tests with the race detector:

```bash
go test -race ./...
```
