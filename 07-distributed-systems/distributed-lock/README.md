# Distributed Lock

## Problem

Multiple service processes may attempt to execute the same critical operation at
the same time. A process-local mutex cannot coordinate work across machines or
application instances.

## Design

Redis acts as a shared lock coordinator.

```text
Instance A ─┐
Instance B ─┼── Redis lock ── protected resource
Instance C ─┘
```

Acquisition uses:

```text
SET lock:<resource> <token> NX PX <ttl>
```

Each lease has a unique ownership token and a finite TTL.

## Implementation

The project uses Go and Redis.

- `internal/lock` implements acquisition, extension, release, and keepalive.
- `internal/critical` executes work while holding a lease.
- `internal/api` exposes the critical-section endpoint.
- `cmd/server` starts independently configurable service instances.

Release and extension use atomic Lua scripts that modify a lock only when its
token still matches the current owner.

## Failure Cases

- Acquisition times out when another process owns the resource.
- TTL releases a lock after an owner crashes.
- Keepalive extends leases for long-running work.
- A stale owner cannot release a newer owner's lock.
- Redis failure prevents lock acquisition.
- A single Redis instance is a coordination dependency and failure point.
- Long process pauses can allow a lease to expire unexpectedly.
- Stronger systems may require fencing tokens to reject stale owners.

## What I Learned

- A local mutex cannot coordinate separate processes.
- Distributed locks require unique ownership tokens.
- Lock acquisition must be atomic.
- Release must compare ownership before deleting.
- TTL prevents permanent locks after crashes.
- Lease renewal protects long-running critical sections.
- Distributed locks reduce concurrency but do not replace idempotency.

## Running

Start Redis:

```bash
docker compose up -d
```

Run instance one:

```bash
INSTANCE_ID=instance-1 \
SERVER_ADDRESS=:8088 \
go run ./cmd/server
```

Run instance two:

```bash
INSTANCE_ID=instance-2 \
SERVER_ADDRESS=:8089 \
go run ./cmd/server
```

Execute a protected operation:

```bash
curl -X POST http://localhost:8088/api/critical/daily-report \
  -H "Content-Type: application/json" \
  -d '{"workDurationMs":5000}'
```

Run tests:

```bash
go test ./...
```

Stop Redis:

```bash
docker compose down
```