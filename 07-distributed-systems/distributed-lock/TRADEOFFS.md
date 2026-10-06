# Trade-offs

## Redis coordinator

Atomic Redis commands simplify acquisition and ownership checks, but Redis is a
coordination dependency whose failover semantics affect correctness.

## Leases

TTL releases abandoned locks, while keepalive supports longer work. Process
pauses or network delays can still let a lease expire while its old owner runs.

## Missing fencing tokens

Compare-and-delete prevents stale release, but not stale writes to the protected
resource. A monotonically increasing fencing token is the preferred extension.
