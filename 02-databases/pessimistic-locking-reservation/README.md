# Pessimistic Locking Reservation

## Problem

Concurrent requests can reserve the same limited resource and oversell its capacity if they read and update it independently.

```text
Request A reads available=1
Request B reads available=1
Both create a reservation
```

The system must serialize access to each resource before validating and reducing capacity.

## Design

```text
Client → Controller → Service → Locked Repository → PostgreSQL
```

The repository obtains a pessimistic write lock:

```java
@Lock(LockModeType.PESSIMISTIC_WRITE)
```

PostgreSQL executes a query similar to:

```sql
SELECT *
FROM reservation_resources
WHERE id = ?
FOR UPDATE;
```

A concurrent transaction requesting the same row waits until the current transaction commits or rolls back.

## Implementation

- Java 21 and Spring Boot
- Spring Data JPA and PostgreSQL
- Flyway-managed schema
- Transactional reservation service
- Pessimistic write locking
- Capacity validation and database constraints
- JUnit and Testcontainers

The locked transaction reads the resource, validates capacity, creates a reservation, updates available capacity, and commits atomically.

## Failure Cases

- Invalid request → `400 Bad Request`
- Missing resource → `404 Not Found`
- Duplicate resource code → `409 Conflict`
- Insufficient capacity → `409 Conflict`
- Lock acquisition failure → `503 Service Unavailable`

A rollback removes partial changes and releases any held database locks.

## What I Learned

- How pessimistic row locking serializes access.
- How `SELECT FOR UPDATE` behaves in PostgreSQL.
- Why locks require an active transaction.
- How locks apply per row rather than to the whole table.
- How to test waiting behavior with threads and latches.
- When pessimistic locking is preferable to optimistic locking.

## Running

Start PostgreSQL:

```bash
docker compose up -d
```

Run the application:

```bash
mvn spring-boot:run
```

Run all tests:

```bash
mvn test
```

Stop PostgreSQL:

```bash
docker compose down
```