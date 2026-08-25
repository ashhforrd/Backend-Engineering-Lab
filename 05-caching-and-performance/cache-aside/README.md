# Cache Aside

## Problem

Repeatedly reading the same data from PostgreSQL increases database load and
response latency. Redis can serve frequently accessed data faster, but it must
not become the source of truth.

## Design

The application manages the cache explicitly.

```text
Read:
Redis hit  -> return cached product
Redis miss -> PostgreSQL -> populate Redis -> return product

Write:
update PostgreSQL -> delete cached product
```

PostgreSQL remains the source of truth. Redis stores temporary copies with a
30-second TTL.

## Implementation

The project uses Go, PostgreSQL, and Redis.

- `internal/product` contains the cache-aside business flow.
- `internal/postgres` implements the product repository.
- `internal/redis` implements the product cache.
- `internal/api` exposes read and update endpoints.
- `cmd/server` wires the application.
- `migrations` creates and seeds the products table.

Responses include `CACHE` or `DATABASE` to make the read source visible.

## Failure Cases

- A cache miss falls back to PostgreSQL.
- A Redis failure also falls back to PostgreSQL.
- A failed cache write does not fail a successful database read.
- A failed invalidation can serve stale data until the TTL expires.
- Concurrent cache misses can cause a cache stampede.
- Updating the database before invalidating cache avoids deleting the cache
  before a failed database update.

## What I Learned

- Cache-aside keeps cache management in the application.
- Cache is an optimization, not the source of truth.
- TTL limits how long stale data can survive.
- Writes require explicit cache invalidation.
- Cache failures should degrade performance, not availability.
- Interfaces separate business logic from PostgreSQL and Redis.

## Running

Start PostgreSQL and Redis:

```bash
docker compose up -d
```

Run the application:

```bash
go run ./cmd/server
```

Read a product:

```bash
curl http://localhost:8085/api/products/1
```

Update a product:

```bash
curl -X PUT http://localhost:8085/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard Pro",
    "price": 1500000,
    "stock": 8
  }'
```

Run tests:

```bash
go test ./...
```

Stop local infrastructure:

```bash
docker compose down
```