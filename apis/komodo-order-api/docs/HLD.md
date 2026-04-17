# Order API — High-Level Design

## Order Number Schema

### Decision
Order entities use two separate identifiers:

1. **Internal ID** (`id`) — machine identity, used as the DynamoDB primary key.
   Format: `{TYPE}-{zero-padded 6-digit sequence}`
   Examples: `ORD-001234`, `RTN-000089`, `EXC-000017`

2. **Display ID** (`displayId`) — customer-facing label, rendered at write time and stored as a plain attribute. Never used as a key or queried against.
   - Root order: `001234`
   - First return on that order: `001234-R1`
   - Second return: `001234-R2`
   - First exchange: `001234-X1`

### Rationale
Encoding the relationship in the display ID (e.g. `001234-R1`) gives customers an immediately readable link between a return/exchange and its originating order without requiring a lookup. The postfix format was chosen over a prefix (`RTN-001234`) because the order number is the primary identity — the type qualifier is secondary.

The internal ID and display ID are intentionally decoupled:
- The internal ID is always unique and type-safe for DB operations
- The display ID is a formatting convention that encodes the parent relationship visually
- Returns and exchanges have their own independent sequences, so they never consume slots in the main order counter

### Derivative entities (Returns, Exchanges)
Returns and exchanges are independent entities with their own internal IDs and sequences. The relationship to the parent order is carried by an explicit `parentOrderId` FK — not encoded in the internal ID.

This means:
- A return is always addressable by its own ID (`RTN-000089`) for internal operations
- The display ID (`001234-R1`) is derived from `parentOrderId` + a 1-based child index at render time
- Lookups by display ID strip the suffix to find the parent order, then resolve the child
