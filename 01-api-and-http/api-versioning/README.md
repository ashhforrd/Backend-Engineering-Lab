# API Versioning

## Problem

An API contract must evolve without suddenly breaking clients that still depend
on its previous response structure.

## Design

The server exposes two URL-based API versions:

- `GET /api/v1/users/{id}` preserves the legacy `name` field.
- `GET /api/v2/users/{id}` introduces `profile`, `status`, and `createdAt`.

Both versions share the same domain model, service, and repository. Each handler
maps the domain user into its own version-specific response DTO.

V1 remains operational while advertising its lifecycle through the standardized
`Deprecation` and `Sunset` headers and a link to its successor version.

## Implementation

- Go standard library HTTP server and method-aware routing
- Separate handlers and response DTOs for v1 and v2
- Shared user service and in-memory repository
- Consistent JSON error responses
- Contract tests using `httptest`

## Failure Cases

- Unknown users return `404 Not Found` in both versions.
- Empty user IDs are rejected by the shared service.
- Repository and unexpected failures become `500 Internal Server Error`.
- V1 response changes are detected by contract tests.

## What I Learned

- API versions should separate transport contracts from business logic.
- A new response shape can coexist with an old one during client migration.
- Deprecation is a lifecycle signal; it does not immediately disable an API.
- Contract tests protect backward compatibility during future refactoring.

## Running

Start the server:

```bash
go run ./cmd/server
```

Call both versions:

```bash
curl -i http://localhost:8080/api/v1/users/user-001
curl -i http://localhost:8080/api/v2/users/user-001
```

Run the tests:

```bash
go test ./...
```

References: [RFC 9745](https://www.rfc-editor.org/rfc/rfc9745.html) and
[RFC 8594](https://www.rfc-editor.org/rfc/rfc8594.html).
