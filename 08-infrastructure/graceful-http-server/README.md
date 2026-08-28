# Graceful HTTP Server

## Problem

Containers and orchestrators stop applications with `SIGTERM`. Exiting
immediately can terminate active requests and background work, causing client
errors and incomplete operations.

## Design

```text
SIGTERM
  ↓
mark instance NOT_READY
  ↓
wait for load-balancer propagation
  ↓
stop accepting new HTTP traffic
  ↓
drain active requests
  ↓
cancel background work
  ↓
exit process
```

A shutdown deadline prevents the process from waiting forever.

## Implementation

The project uses the Go standard library.

- `internal/lifecycle` tracks readiness, active requests, and shutdown order.
- `internal/api` exposes liveness, readiness, and slow-work endpoints.
- `internal/background` demonstrates context-aware background work.
- `cmd/server` handles `SIGINT`, `SIGTERM`, server errors, and goroutine cleanup.

Health endpoints:

```text
GET /health/live
GET /health/ready
```

Slow request simulator:

```text
GET /api/work?durationMs=4000
```

## Failure Cases

- New requests receive `503` after shutdown begins.
- Existing HTTP handlers are allowed to finish.
- Requests exceeding the shutdown deadline are forcefully closed.
- Background goroutines require explicit cancellation and waiting.
- `http.Server.Shutdown` does not automatically wait for arbitrary goroutines.
- A propagation delay is needed because load-balancer updates are not instant.

## What I Learned

- `SIGTERM` requests process termination.
- Graceful shutdown still ends the process.
- Readiness controls traffic while liveness indicates process health.
- HTTP draining waits for active handlers to return.
- Context cancellation stops background components.
- Wait groups prevent main from exiting before goroutines finish.
- Forced shutdown provides an upper bound on termination time.

## Running

Build the server:

```bash
go build -o /tmp/graceful-http-server ./cmd/server
```

Run it:

```bash
/tmp/graceful-http-server
```

Start a slow request:

```bash
curl "http://localhost:8088/api/work?durationMs=4000"
```

Send `SIGTERM`:

```bash
SERVER_PID=$(pgrep -f '^/tmp/graceful-http-server$')
kill -TERM "$SERVER_PID"
```

Run tests:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```