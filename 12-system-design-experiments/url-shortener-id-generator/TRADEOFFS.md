# Trade-offs

## Snowflake-style IDs

Generation is local and fast, but operations must assign unique node IDs. Values
also reveal approximate creation time and ordering.

## Bit allocation

More node bits support more instances but reduce time range or per-millisecond
throughput. This layout supports 1,024 nodes and 4,096 IDs per node per ms.

## Clock and encoding

Clock rollback is rejected instead of risking duplicates. Base62 shortens the
number but is reversible and provides no secrecy or authorization.
