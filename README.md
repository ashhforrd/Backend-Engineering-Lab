# Backend Engineering Lab

Hands-on implementations of backend and systems-engineering concepts. Each lab
isolates one production problem, makes the failure mode observable, and records
the design decisions and trade-offs.

## Language Strategy

- **Go** — HTTP services, reliability, infrastructure, observability, and async workflows.
- **Java** — transactions, database behavior, persistence, and enterprise service patterns.
- **Rust** — security-sensitive code, memory safety, and high-integrity system components.
- **C++** — low-level concurrency, synchronization, memory, and performance.
- **Python** — reserved for DSA practice and coding challenges, outside these backend labs.

The goal is not to use every language for every problem. Each language is chosen
when its model makes the engineering concept easier to observe.

## Labs

| Chapter | Project | Language | Main concept |
|---|---|---|---|
| API & HTTP | [Idempotency Key API](01-api-and-http/idempotency-key-api) | Go | Safe request retries |
| API & HTTP | [API Versioning](01-api-and-http/api-versioning) | Go | Backward-compatible contracts |
| API & HTTP | [Rate Limit Headers](01-api-and-http/api-rate-limit-headers) | Go | Quota communication and `429` |
| Databases | [Optimistic Locking Inventory](02-databases/optimistic-locking-inventory) | Java | Version-based conflict detection |
| Databases | [Pessimistic Locking Reservation](02-databases/pessimistic-locking-reservation) | Java | Serialized row access |
| Databases | [Transaction Isolation Lab](02-databases/transaction-isolation-lab) | Java | Isolation anomalies and snapshots |
| Databases | [Database Index Benchmark](02-databases/db-index-benchmark) | Java | Query plans and index selectivity |
| Concurrency | [Bounded Worker Pool](03-concurrency/bounded-worker-pool) | Go | Controlled concurrency and backpressure |
| Concurrency | [Producer Consumer](03-concurrency/producer-consumer) | C++ | Mutexes and condition variables |
| Reliability | [Retry with Exponential Backoff](04-reliability/retry-with-exponential-backoff) | Go | Retry policy, jitter, and cancellation |
| Reliability | [Circuit Breaker](04-reliability/circuit-breaker) | Go | Dependency failure isolation |
| Caching | [Cache Aside](05-caching-and-performance/cache-aside) | Go | Cache consistency and invalidation |
| Async Processing | [Basic Job Queue](06-async-processing/basic-job-queue) | Go | Producer/consumer job execution |
| Distributed Systems | [Distributed Lock](07-distributed-systems/distributed-lock) | Go | Cross-process coordination and leases |
| Infrastructure | [Graceful HTTP Server](08-infrastructure/graceful-http-server) | Go | Draining and signal handling |
| Observability | [Structured Logging](09-observability/structured-logging) | Go | Machine-searchable contextual logs |
| Testing | [Table-Driven Tests](10-testing/table-driven-tests) | Go | Structured test cases and boundaries |
| Security | [Password Hashing](11-security/password-hashing) | Rust | Argon2id credential storage |
| System Design | [URL Shortener ID Generator](12-system-design-experiments/url-shortener-id-generator) | Rust | Distributed unique ID generation |

## Repository Principles

- Reproduce a real failure before implementing the protection.
- Keep business logic separate from transport and infrastructure concerns.
- Test concurrency and failure paths, not only happy paths.
- Document what the design does **not** guarantee.
- Prefer measurable evidence: execution plans, counters, timings, and race checks.

## Running a Lab

Each project is self-contained and has its own README. Common commands:

```bash
# Go
go test ./...

# Java
mvn test

# Rust
cargo test

# C++
cmake -S . -B build
cmake --build build
ctest --test-dir build --output-on-failure
```

Projects using PostgreSQL or Redis include a local `docker-compose.yml`.

## Continuous Integration

GitHub Actions builds and tests every Go, Java, Rust, and C++ lab independently.
This keeps the polyglot structure intentional while still enforcing one quality
baseline across the repository.

## Suggested Learning Order

Start with API semantics, continue through database consistency and local
concurrency, then move into failure handling, caching, asynchronous processing,
distributed coordination, infrastructure, observability, testing, and security.
The system-design experiments combine concepts from the earlier chapters.
