# Bounded Worker Pool

## Problem

Starting one goroutine for every task can overload memory, databases, file
descriptors, or downstream services. Concurrency needs explicit limits so the
system remains stable under load.

## Design

```text
Producer
   ↓
Bounded task channel
   ↓
Fixed worker goroutines
   ↓
Result channel
   ↓
Result collector
```

The worker count limits active tasks. Queue capacity limits waiting tasks.
Submissions are rejected when the queue is full.

## Implementation

The project uses the Go standard library.

- `internal/pool` implements task submission, workers, results, and metrics.
- `internal/api` generates workloads and exposes pool statistics.
- `cmd/server` starts the pool, result collector, and HTTP server.

Each task contains an ID and a function:

```go
type Task struct {
    ID  string
    Run TaskFunc
}
```

Atomic counters track active, completed, failed, and maximum observed workers.

## Failure Cases

- Invalid configuration prevents pool creation.
- Submission before startup returns `ErrNotRunning`.
- A full queue returns `ErrQueueFull`.
- Task errors are returned through the result channel.
- Task panics are recovered and converted into errors.
- Shutdown cancels active tasks instead of draining all queued work.
- Workers can block if no consumer drains the result channel.

## What I Learned

- Worker count controls concurrency.
- Queue capacity creates bounded backpressure.
- Channels distribute tasks safely between goroutines.
- Atomic operations protect counters without a mutex.
- Compare-and-swap safely tracks a concurrent maximum.
- Panic recovery isolates failures from the worker process.
- Controlled throughput is often safer than maximum throughput.

## Running

Start the application:

```bash
go run ./cmd/server
```

Submit a workload:

```bash
curl -X POST http://localhost:8087/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "taskCount": 20,
    "durationMs": 2000,
    "failEvery": 5
  }'
```

Inspect statistics:

```bash
curl http://localhost:8087/api/stats
```

Run tests:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```