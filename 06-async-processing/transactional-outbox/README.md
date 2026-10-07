# Transactional Outbox

## Problem

Saving an order in PostgreSQL and publishing an event to Kafka are independent
operations. A crash between them can leave committed data without an event.

## Design

The order and outbox event commit in one PostgreSQL transaction:

~~~text
POST /api/orders
        |
        v
  BEGIN TRANSACTION
    INSERT orders
    INSERT outbox_events
  COMMIT
        |
        v
 scheduled publisher -> Kafka -> idempotent consumer
~~~

The database write is atomic. Kafka delivery is asynchronous and at least once,
so consumers deduplicate messages using the event ID.

## Implementation

- Java 21 and Spring Boot
- PostgreSQL with Flyway
- Kafka in KRaft mode
- FOR UPDATE SKIP LOCKED for concurrent publishers
- Partial index for pending events
- Bounded retries and terminal FAILED events
- Manual recovery endpoint
- Consumer deduplication through processed_events

Outbox lifecycle:

~~~text
PENDING --success--> PUBLISHED
   |
   +--retry limit--> FAILED --manual retry--> PENDING
~~~

## Failure Cases

- Kafka unavailable: the order commits and its event remains pending.
- Publisher crashes before send: a later run picks the event up.
- Publisher crashes after send but before commit: Kafka can receive a duplicate.
- Repeated failures: the event becomes failed after five attempts.
- Duplicate delivery: the consumer ignores an event ID already processed.

## What I Learned

- One local transaction cannot atomically commit PostgreSQL and Kafka.
- The outbox changes cross-system delivery into an atomic database write and a
  retryable relay.
- At-least-once delivery requires idempotent consumers.
- Kafka partitions preserve ordering per key and enable parallel consumption.
- SKIP LOCKED lets multiple publishers safely share pending work.

## Running

Start infrastructure:

~~~bash
docker compose up -d
~~~

Run the application:

~~~bash
mvn spring-boot:run
~~~

Create an order:

~~~bash
curl -s -X POST http://localhost:8082/api/orders \
  -H "Content-Type: application/json" \
  -d '{"customerId":"customer-001","totalAmount":250000}'
~~~

Inspect its outbox events:

~~~bash
curl -s http://localhost:8082/api/orders/ORDER_ID/events
~~~

Replay a failed event:

~~~bash
curl -s -X POST http://localhost:8082/api/outbox/EVENT_ID/retry
~~~

Run tests:

~~~bash
mvn test
~~~
