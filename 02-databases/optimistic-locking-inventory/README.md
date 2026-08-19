# Optimistic Locking Inventory

## Overview

A Spring Boot inventory API demonstrating version-based conflict detection with JPA and PostgreSQL.

## Problem

Concurrent requests can read the same stock and overwrite each other, causing lost updates.

```text
A reads 10 ─┐
B reads 10 ─┴─ both write 9 → one update is lost
```

## Requirements

- Create and retrieve inventory items.
- Decrease stock without allowing negative quantities.
- Reject duplicate SKUs and invalid requests.
- Detect concurrent updates and return `409 Conflict`.

## Architecture

```text
Client → Controller → Service → Repository → Hibernate → PostgreSQL
```

The service uses Java 21, Spring Boot, Spring Data JPA, Flyway, Docker, and Testcontainers.

## Core Design Decisions

- `@Version` enables optimistic locking.
- `@Transactional` defines atomic operations.
- DTOs separate the HTTP contract from JPA entities.
- Entity methods enforce inventory rules.
- Flyway manages the database schema.

## Data Model

`inventory_items` stores `id`, `sku`, `name`, `quantity`, `version`, `created_at`, and `updated_at`.

PostgreSQL enforces unique SKUs and non-negative quantities.

## Reliability

Validation, entity rules, database constraints, and transactions protect inventory integrity at different layers.

## Concurrency

Hibernate generates version-aware updates:

```sql
UPDATE inventory_items
SET quantity = ?, version = ?
WHERE id = ? AND version = ?;
```

Only one transaction can update a given version. A transaction using an outdated version is rolled back.

## Failure Scenarios

- Invalid request → `400 Bad Request`
- Missing item → `404 Not Found`
- Duplicate SKU → `409 Conflict`
- Insufficient inventory → `409 Conflict`
- Concurrent update → `409 Conflict`

## Observability

Hibernate SQL logging exposes the generated version-aware queries. API errors include status, message, path, timestamp, and field errors.

## Benchmarks

No formal performance benchmark has been performed. The integration test verifies that exactly one of two concurrent updates succeeds.

## Trade-offs

Optimistic locking avoids read locks and works well under low contention. Under high contention, more requests fail and may require bounded retries.

## Testing

Run:

```bash
mvn test
```

Unit tests verify inventory rules. A Testcontainers integration test uses two synchronized transactions against real PostgreSQL.

## Running Locally

```bash
docker compose up -d
mvn spring-boot:run
```

The API runs at `http://localhost:8080`.

Stop the infrastructure with:

```bash
docker compose down
```

## What I Learned

- How lost updates occur.
- How JPA `@Version` detects conflicts.
- How transactions and dirty checking work.
- How to test concurrent database operations.

## Future Work

- Add bounded retry policies.
- Add metrics and health endpoints.
- Compare optimistic and pessimistic locking.
- Add API-level integration tests and benchmarks.