# Benchmark Methodology

A single timing can be distorted by cache state, background activity, and query
planning overhead. Use repeated runs and report a range rather than presenting
one local measurement as universal.

Start PostgreSQL, then run five isolated dataset builds:

```bash
docker compose up -d
./scripts/run-benchmark.sh 5
```

The script reports minimum, average, and maximum execution time for indexed and
unindexed queries. Each run recreates the one-million-row table so the phases
remain comparable. Record the machine, PostgreSQL version, dataset size, query,
execution plan, buffer activity, and run count alongside the timings.

The existing README result is an observed example, not a performance guarantee.
