# Transaction Isolation Lab

## Problem

Concurrent transactions can observe inconsistent data depending on their isolation level.

This lab demonstrates:

- dirty reads;
- non-repeatable reads;
- phantom reads;
- PostgreSQL isolation behavior.

## Design

```text
HTTP Request
    ↓
Experiment Controller
    ↓
Two JDBC Connections
    ↓
Two Concurrent Transactions
    ↓
PostgreSQL
```

Each experiment controls `Connection`, isolation level, commit, and rollback manually.

## Implementation

The supported levels are:

```text
READ_UNCOMMITTED
READ_COMMITTED
REPEATABLE_READ
SERIALIZABLE
```

PostgreSQL treats `READ UNCOMMITTED` as `READ COMMITTED`, so dirty reads are prevented.

PostgreSQL may report the configured level as `read uncommitted`, but its visibility behavior still prevents dirty reads.

Expected behavior:

| Level | Dirty read | Non-repeatable read | Phantom read |
|---|---:|---:|---:|
| READ UNCOMMITTED | Prevented | Possible | Possible |
| READ COMMITTED | Prevented | Possible | Possible |
| REPEATABLE READ | Prevented | Prevented | Prevented |
| SERIALIZABLE | Prevented | Prevented | Prevented |

The project uses Java 21, Spring Boot JDBC, PostgreSQL, Flyway, JUnit, and Testcontainers.

## Failure Cases

- Unsupported isolation value → `400 Bad Request`
- PostgreSQL unavailable → experiment fails
- Missing seed data → experiment fails
- Insufficient connection pool → concurrent transaction cannot start
- Concurrent experiment requests may interfere with shared reset data

This lab is intended to run one experiment at a time.

## What I Learned

- How JDBC configures transaction isolation.
- Why PostgreSQL does not allow dirty reads.
- How `READ COMMITTED` creates a new snapshot per statement.
- How `REPEATABLE READ` keeps one transaction snapshot.
- The difference between changed rows and phantom rows.
- How commit timing affects transaction visibility.
- How to test isolation behavior with real PostgreSQL.

## Running

Start PostgreSQL:

```bash
docker compose up -d
```

Run the application:

```bash
mvn spring-boot:run
```

The API runs at `http://localhost:8082`.

Examples:

```bash
curl -X POST \
  "http://localhost:8082/api/experiments/dirty-read?isolation=READ_UNCOMMITTED"

curl -X POST \
  "http://localhost:8082/api/experiments/non-repeatable-read?isolation=READ_COMMITTED"

curl -X POST \
  "http://localhost:8082/api/experiments/phantom-read?isolation=REPEATABLE_READ"
```

Run tests:

```bash
mvn test
```

Stop PostgreSQL:

```bash
docker compose down
```
