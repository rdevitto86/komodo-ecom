# Order API — Data Model

## Tables

### Orders
Primary table for root purchase entities.

| Attribute     | Type   | Notes |
|---------------|--------|-------|
| `PK`          | String | `ORD-{sequence}` — internal ID |
| `displayId`   | String | `{sequence}` zero-padded, e.g. `001234` |
| `userId`      | String | FK → user-api |
| `status`      | String | See order status lifecycle |
| `items`       | List   | Embedded `OrderItem` records |
| `address`     | Map    | Shipping address snapshot |
| `payment`     | Map    | Payment method + transaction ref |
| `totals`      | Map    | Subtotal, tax, shipping, discount, total |
| `createdAt`   | String | ISO 8601 |
| `updatedAt`   | String | ISO 8601 |

**GSI: UserOrdersIndex**
- PK: `userId`, SK: `createdAt`
- Access pattern: "all orders for a user, newest first"

---

### Returns
Independent table for return entities. Linked to Orders by FK.

| Attribute       | Type   | Notes |
|-----------------|--------|-------|
| `PK`            | String | `RTN-{sequence}` — internal ID |
| `displayId`     | String | `{parentSeq}-R{n}`, e.g. `001234-R1` |
| `parentOrderId` | String | FK → Orders PK |
| `userId`        | String | FK → user-api |
| `status`        | String | See order status lifecycle |
| `items`         | List   | Embedded `OrderItem` records (subset of parent) |
| `reason`        | String | Customer-provided return reason |
| `createdAt`     | String | ISO 8601 |
| `updatedAt`     | String | ISO 8601 |

**GSI: ReturnsByOrderIndex**
- PK: `parentOrderId`, SK: `createdAt`
- Access pattern: "all returns for a given order"

---

### Exchanges
Independent table for exchange entities. Linked to Orders by FK.

| Attribute       | Type   | Notes |
|-----------------|--------|-------|
| `PK`            | String | `EXC-{sequence}` — internal ID |
| `displayId`     | String | `{parentSeq}-X{n}`, e.g. `001234-X1` |
| `parentOrderId` | String | FK → Orders PK |
| `userId`        | String | FK → user-api |
| `status`        | String | See order status lifecycle |
| `returnItems`   | List   | Items being returned |
| `exchangeItems` | List   | Replacement items |
| `priceDelta`    | Number | Positive = customer owes, negative = customer credited |
| `createdAt`     | String | ISO 8601 |
| `updatedAt`     | String | ISO 8601 |

**GSI: ExchangesByOrderIndex**
- PK: `parentOrderId`, SK: `createdAt`
- Access pattern: "all exchanges for a given order"

---

## Sequence Counters (Redis)

| Key            | Tracks          |
|----------------|-----------------|
| `seq:order`    | Next ORD number |
| `seq:return`   | Next RTN number |
| `seq:exchange` | Next EXC number |

Incremented atomically via `INCR` before each entity write. Redis must run with `appendonly yes` in non-local environments.

---

## Order Status Lifecycle

```
pending → confirmed → processing → shipped → delivered
                                           ↘ cancelled
                              ← refunded ←
```

Returns and exchanges follow the same status field with context-appropriate transitions.
