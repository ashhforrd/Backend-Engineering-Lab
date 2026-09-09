# Database Index Benchmark

## Problem

Queries over large tables become expensive when PostgreSQL must inspect nearly
every row to find a small result set.

## Design

The benchmark creates one million deterministic orders and runs the same query
before and after adding this composite index:

```sql
CREATE INDEX idx_orders_customer_created_at
    ON orders (customer_id, created_at DESC);
```

The index matches both the equality filter and requested ordering:

```sql
WHERE customer_id = ?
ORDER BY created_at DESC
LIMIT 20
```

## Implementation

- Java 21 and JDBC
- PostgreSQL 17 in Docker Compose
- SQL dataset generation with `generate_series`
- `EXPLAIN (ANALYZE, BUFFERS)` for measured execution plans
- Warm-up execution before each measured query
- Environment-based database configuration

## Results

One local run produced:

| Metric | Without index | With index |
|---|---:|---:|
| Execution time | 41.004 ms | 0.156 ms |
| Shared buffer hits | 9,334 | 13 |
| Main scan | Parallel Seq Scan | Bitmap Index Scan |

The indexed query was about 263 times faster in this run. Results vary with
hardware, cache state, data distribution, and PostgreSQL planner decisions.

## Failure Cases

- The benchmark cannot start while PostgreSQL is unavailable on port `5435`.
- Reusing an old dataset can distort results; `setup.sql` recreates the table.
- Low-selectivity queries may still favor a sequential scan.
- Indexes consume storage and add work to inserts, updates, and deletes.
- The simple SQL runner is intended for these controlled scripts and does not
  parse semicolons inside SQL strings or procedural blocks.

## What I Learned

- Index usefulness depends on query shape and selectivity.
- Execution plans explain performance better than timing alone.
- `ANALYZE` supplies table statistics to the query planner.
- A composite index can support filtering and ordering together.
- PostgreSQL may choose a bitmap scan instead of a plain index scan.

## Running

```bash
docker compose up -d
mvn test
mvn compile exec:java
```

Stop PostgreSQL:

```bash
docker compose down
```
