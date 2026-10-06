# URL Shortener ID Generator

## Problem

A URL shortener needs compact identifiers that remain unique while many service
instances generate IDs concurrently without coordinating on every request.

## Design

The generator builds a Snowflake-style 63-bit value:

```text
| timestamp: 41 bits | node ID: 10 bits | sequence: 12 bits |
```

The numeric value is Base62-encoded for a shorter URL-safe representation.

## Implementation

- Rust with a generic clock abstraction
- 1,024 possible node IDs
- 4,096 IDs per node per millisecond
- Mutex-protected timestamp and sequence state
- Clock rollback detection
- Sequence exhaustion waits for the next millisecond
- Concurrent uniqueness and Base62 tests

## Failure Cases

- Duplicate node IDs can generate collisions across processes.
- A clock moving backward is rejected rather than risking duplicate IDs.
- Sequence exhaustion temporarily waits for time to advance.
- The 41-bit timestamp eventually exhausts its representable range.
- IDs expose approximate creation time and node information.

## What I Learned

- Distributed uniqueness can be composed from time, node identity, and sequence.
- Bit allocation determines capacity and system lifetime.
- A short representation is encoding, not encryption.
- Operational node-ID assignment is part of correctness.

## Running

```bash
cargo run
cargo test
```
