# Order API — Low-Level Design

## Sequence counter
Sequence numbers are generated via **Redis `INCR`** — one key per type (`seq:order`, `seq:return`, `seq:exchange`). This gives atomic, sub-millisecond increments with no DynamoDB hot-partition risk.

Redis must be run with persistence enabled (`appendonly yes`) in all non-local environments. Losing the counter would produce duplicate display IDs.

## CS tooling
The lookup tool strips any suffix/prefix and resolves by numeric sequence. An agent entering `001234`, `ORD-001234`, or `001234-R1` all resolve to the same order. The prefix/suffix is display metadata, not a query requirement.
