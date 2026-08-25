# Circuit Breaker

## Problem

A continuously failing dependency can consume connections, increase latency, and
cause cascading failures. Calling it repeatedly wastes resources when it has
little chance of succeeding.

## Design

The circuit breaker has three states:

```text
CLOSED -> OPEN -> HALF_OPEN
   ^                   |
   +-------------------+
```

- `CLOSED`: requests reach the dependency.
- `OPEN`: requests fail immediately.
- `HALF_OPEN`: limited probe requests test whether the dependency recovered.

## Implementation

The project uses Go's standard library.

- `internal/circuitbreaker` manages state transitions and request admission.
- `internal/downstream` calls the product dependency.
- `internal/product` applies the breaker to the business flow.
- `internal/simulator` simulates a healthy or failing dependency.
- `internal/api` exposes product and breaker-status endpoints.
- `cmd/server` wires dependencies and starts the server.

State is protected by a mutex because concurrent requests share one breaker.
Generation numbers prevent old requests from modifying a newer breaker state.

## Failure Cases

- Three consecutive dependency failures open the circuit.
- Requests fail fast while the circuit is open.
- After the open timeout, one probe enters the half-open state.
- A successful probe closes the circuit.
- A failed probe opens it again.
- Client errors do not count as dependency failures.

## What I Learned

- Circuit breakers prevent cascading failures.
- Closed circuits allow traffic; open circuits block it.
- Shared state must be concurrency-safe.
- Not every application error means a dependency is unhealthy.
- Half-open probes test recovery without restoring full traffic.
- Time and state transitions should be testable independently.

## Running

Start the server:

```bash
go run ./cmd/server
```

Read a product:

```bash
curl http://localhost:8084/api/products/product-001
```

Disable the simulated dependency:

```bash
curl -X PUT http://localhost:8084/simulator/availability \
  -H "Content-Type: application/json" \
  -d '{"available":false}'
```

Inspect the breaker:

```bash
curl http://localhost:8084/api/circuit-breaker
```

Enable the dependency again:

```bash
curl -X PUT http://localhost:8084/simulator/availability \
  -H "Content-Type: application/json" \
  -d '{"available":true}'
```