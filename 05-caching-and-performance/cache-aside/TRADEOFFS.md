# Trade-offs

## Explicit application control

Cache-aside keeps PostgreSQL as the source of truth and makes fallback behavior
clear, but every call site must follow the same read and invalidation rules.

## Availability over freshness

Redis failures fall back to PostgreSQL. This preserves availability but can
overload the database during a broad cache outage.

## Invalidation and stampedes

Database-first updates avoid cache values newer than the database, but failed
invalidation can serve stale data until TTL expiry. Concurrent misses can also
stampede PostgreSQL; singleflight and TTL jitter are future extensions.
