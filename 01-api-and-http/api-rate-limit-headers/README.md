# API Rate Limit Headers

## Problem

Rejecting excessive traffic is not enough. Clients also need to know their
quota, how much remains, and when they can safely retry.

## Design

The service uses a fixed window per API key, with the remote IP address as a
fallback identity. A mutex protects the in-memory counters from concurrent
requests.

Every protected response includes:

- `X-RateLimit-Limit`: maximum requests in the window
- `X-RateLimit-Remaining`: quota left after the current request
- `X-RateLimit-Reset`: window reset time as a Unix timestamp
- `Retry-After`: seconds to wait after receiving `429`

## Implementation

- Go standard library HTTP server
- Configurable request limit and window duration
- Middleware that protects selected routes
- Independent counters for each client
- Clock abstraction for deterministic tests
- Unrestricted health endpoint

## Failure Cases

- Invalid limiter configuration prevents startup.
- Exhausted quota returns `429 Too Many Requests`.
- Missing API keys fall back to the client IP address.
- In-memory quota is lost when the process restarts.

## What I Learned

- Rate limiting and communicating rate limits are separate concerns.
- Middleware can enforce a policy without changing endpoint logic.
- `Retry-After` lets clients implement responsible retry behavior.
- Shared mutable counters require synchronization.

## Running

Start the server:

```bash
go run ./cmd/server
```

Send six requests with the same client identity:

```bash
for i in {1..6}; do
  curl -i -H "X-API-Key: client-001" http://localhost:8081/api/message
done
```

Run the tests:

```bash
go test ./...
```

`429 Too Many Requests` is defined by
[RFC 6585](https://www.rfc-editor.org/rfc/rfc6585.html), and `Retry-After` by
[RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html#name-retry-after).
