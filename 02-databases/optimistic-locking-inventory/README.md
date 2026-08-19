# Optimistic Locking Inventory

## Problem

Concurrent requests can read the same inventory state and overwrite each other:

```text
Request A reads quantity=10
Request B reads quantity=10
Request A writes quantity=9
Request B writes quantity=9
```

Both requests appear successful, but one update is lost.

## Design

```text
Client → Controller → Service → Repository → Hibernate → PostgreSQL
```

The entity uses version-based conflict detection:

```java
@Version
private Long version;
```

Hibernate includes the version in its update condition. A transaction using an outdated version is rejected.

## Implementation

- Java 21 and Spring Boot
- Spring Data JPA and PostgreSQL
- Flyway-managed schema
- `@Transactional` service operations
- DTO validation
- Structured API errors
- JUnit and Testcontainers

Hibernate generates an update similar to:

```sql
UPDATE inventory_items
SET quantity = ?, version = ?
WHERE id = ? AND version = ?;
```

## Failure Cases

- Invalid request → `400 Bad Request`
- Missing item → `404 Not Found`
- Duplicate SKU → `409 Conflict`
- Insufficient inventory → `409 Conflict`
- Concurrent update → `409 Conflict`

## What I Learned

- How lost updates occur.
- How JPA `@Version` detects conflicts.
- How dirty checking and transactions work.
- How to test concurrent updates with `CyclicBarrier`.
- When optimistic locking is appropriate.

## Running

Start PostgreSQL and the application:

```bash
docker compose up -d
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