# Basic Job Queue

## Problem

Slow work such as sending emails or generating reports should not block an HTTP
request. A job queue allows the API to accept work quickly and process it in the
background.

## Design

```text
HTTP Producer
    ↓
Buffered Channel
    ↓
Worker Pool
    ↓
Job Dispatcher
    ↓
Job Handler
```

The API returns `202 Accepted` after storing and enqueueing a job. Clients can
poll the job endpoint to observe its status.

## Implementation

The project uses the Go standard library.

- `internal/producer` creates and enqueues jobs.
- `internal/queue` provides a bounded in-memory queue.
- `internal/worker` runs consumers and dispatches jobs.
- `internal/tasks` implements email and report handlers.
- `internal/job` stores job data and lifecycle state.
- `internal/api` exposes submission and status endpoints.
- `cmd/server` starts the API and three workers.

Job lifecycle:

```text
QUEUED -> PROCESSING -> COMPLETED
                     -> FAILED
```

## Failure Cases

- A full queue rejects new jobs with `503`.
- Invalid job types are rejected with `400`.
- Handler errors mark jobs as `FAILED`.
- Worker cancellation stops queue consumption.
- Jobs and statuses are lost when the process restarts.
- Multiple application instances do not share this in-memory queue.
- There are no retries or durable delivery guarantees.

## What I Learned

- Producers and consumers are decoupled by a queue.
- Buffered channels can implement bounded in-memory queues.
- Worker count controls processing parallelism.
- Backpressure prevents unlimited memory growth.
- Dispatchers route job types to different handlers.
- Shared status data requires concurrency protection.
- `202 Accepted` means work was accepted, not completed.

## Running

Start the application:

```bash
go run ./cmd/server
```

Submit an email job:

```bash
curl -X POST http://localhost:8086/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "SEND_EMAIL",
    "payload": {
      "to": "student@example.com",
      "subject": "Learn Go job queues"
    }
  }'
```

Read job status:

```bash
curl http://localhost:8086/api/jobs/<job-id>
```

Inspect queue capacity:

```bash
curl http://localhost:8086/api/queue
```

Run tests:

```bash
go test ./...
```