# Trade-offs

## In-memory storage

Records disappear on restart and are not shared between instances. A production
design needs a durable shared store with atomic insert-if-absent semantics.

## Global mutex

One mutex makes correctness easy to inspect but serializes unrelated keys.
Per-key coordination or database uniqueness increases concurrency with more
lifecycle complexity.

## Stored response and payload hash

Persisting the response guarantees identical retries, while hashing the payload
detects key reuse with different input. Production records also need TTL, size
limits, canonicalization rules, and a retention policy.
