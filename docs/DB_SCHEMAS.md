# Komodo Platform — Database Schemas

This document catalogs all database schemas across the Komodo platform APIs. Each section groups APIs by their underlying data store technology.

---

## DynamoDB

### komodo-user-api
**Table:** `komodo-users`

```markdown
<!-- Schema template -->
- PK/SK pattern
- GSI definitions
- Attributes
- TTL configuration
```

### komodo-cart-api
**Table:** `komodo-carts`

```markdown
<!-- Schema template -->
- PK/SK pattern (USER#<id> / CART#<id>)
- GSI for guest carts
- Attributes (items, totals, metadata)
- TTL configuration (guest carts only)
```

### komodo-order-api
**Table:** `komodo-orders`

```markdown
<!-- Schema template -->
- PK/SK pattern (ORDER#<id> / LINEITEM#<n>)
- GSI for user order history
- GSI for guest email lookup
- Attributes (status, totals, shipping, etc.)
- Status transition constraints
```

### komodo-payments-api
**Table:** `komodo-payments`

```markdown
<!-- Schema template -->
- PK/SK pattern (CHARGE#<id> / REFUND#<id>)
- GSI for user payment methods
- GSI for payment plans
- Attributes (amount, status, provider metadata)
```

### komodo-shop-inventory-api
**Table:** `komodo-inventory`

```markdown
<!-- Schema template -->
- PK/SK pattern (SKU#<sku> / STOCK#<sku>, HOLD#<holdId>)
- Attributes (available_qty, reserved_qty, committed_qty)
- TTL configuration (hold expiry)
- DynamoDB Streams configuration
```

### komodo-order-reservations-api
**Tables:** `komodo-slots`, `komodo-bookings`

```markdown
<!-- Slots table template -->
- PK/SK pattern (SLOT#<date>_<location>_<type> / SLOT#<id>)
- GSI for date/location filtering
- Attributes (capacity, available, time windows)

<!-- Bookings table template -->
- PK/SK pattern (BOOKING#<id> / METADATA)
- GSI for user booking history
- GSI for guest email lookup
- Attributes (status, service type, repair details)
```

### komodo-shop-promotions-api
**Table:** `komodo-promotions`

```markdown
<!-- Schema template -->
- PK/SK pattern (PROMO#<id> / METADATA)
- GSI for active promotions
- Attributes (discount rules, eligibility, dates)
- First-order tracking fields
```

### komodo-support-api
**Table:** `komodo-chat-sessions`

```markdown
<!-- Schema template -->
- PK/SK pattern (SESSION#<id> / MESSAGE#<n>)
- GSI for user session history
- GSI for anonymous session lookup
- Attributes (messages, status, escalation state)
- TTL configuration (anonymous sessions)
```

### komodo-events-api
**Table:** `komodo-events`

```markdown
<!-- Schema template -->
- PK/SK pattern (EVENT#<id> / METADATA)
- GSI for event type filtering
- GSI for timestamp-based queries
- Attributes (event type, payload, status)
- TTL configuration (event retention)
```

---

## Redis / ElastiCache

### komodo-auth-api
**Purpose:** Token caching, JTI revocation list

```markdown
<!-- Key template -->
- Token cache: `token:<jti>` → token metadata
- Revocation list: `revoked:<jti>` → expiry timestamp
- TTL configuration
```

### komodo-cart-api
**Purpose:** Guest carts, checkout tokens

```markdown
<!-- Key template -->
- Guest cart: `cart:guest:<uuid>` → cart items
- Checkout token: `checkout:<token>` → cart snapshot
- TTL configuration (7d for guest carts)
```

### komodo-order-reservations-api
**Purpose:** Distributed locking for slot booking

```markdown
<!-- Key template -->
- Slot lock: `lock:slot:<slotId>` → lock holder
- TTL configuration (lock expiry)
```

---

## S3

### komodo-shop-items-api
**Bucket:** `komodo-shop-items`

```markdown
<!-- Object key template -->
- Items: `items/{sku}.json`
- Services: `services/{id}.json`
- Inventory: `inventory/{sku}.json`
- Metadata structure
```

### komodo-ssr-engine-svelte
**Bucket:** `komodo-ssr-content`

```markdown
<!-- Object key template -->
- Product pages: `products/{id}/index.html`
- Service pages: `services/{id}/index.html`
- Landing pages: `landing/{slug}/index.html`
- Cache invalidation strategy
```

### komodo-features-api
**Bucket:** `komodo-platform-config`

```markdown
<!-- Object key template -->
- Access config: `{env}/access.json`
- Per-user overrides: `users/{userId}.json`
- Audit log: `audit/{env}/changes.jsonl`
- Cache refresh strategy
```

---

## Typesense

### komodo-search-api
**Collection:** `komodo-catalog`

```markdown
<!-- Schema template -->
- Fields (name, description, category, price, etc.)
- Facet configuration
- Ranking/sorting rules
- Synonym configuration
```

---

## SQLite / RDS

### komodo-statistics-api
**Database:** `komodo-stats` (SQLite → RDS migration planned)

```markdown
<!-- Schema template -->
- Trending items table
- In-cart counts table
- Recently bought table
- Frequently bought with table
- TTL cleanup strategy
```

### komodo-ssr-engine-svelte
**Database:** Local SQLite (content metadata)

```markdown
<!-- Schema template -->
- Content tracking table
- Cache invalidation log
- Local-only storage
```

---

## Shared Cross-Service Considerations

- **Consistent key naming conventions** across DynamoDB tables
- **TTL strategies** for temporary data (guest carts, holds, sessions)
- **GSI patterns** for common query patterns (user history, email lookup)
- **Stream processing** for CDC events (inventory, orders, payments)
- **Backup/restore** procedures per data store
- **Migration paths** (SQLite → RDS, in-memory → DynamoDB)
