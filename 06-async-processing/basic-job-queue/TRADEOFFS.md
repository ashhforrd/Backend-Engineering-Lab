# Trade-offs

## In-memory queue

The design makes producer/consumer behavior visible with little infrastructure,
but accepted jobs disappear if the process crashes.

## Bounded capacity

The queue prevents unlimited memory growth and exposes overload to callers. It
can reject bursts that a durable broker could absorb.

## Delivery semantics

Workers remove jobs before executing them, so a crash can lose work. Durable
systems add acknowledgements, visibility timeouts, retry schedules, DLQs, and
idempotent consumers.
