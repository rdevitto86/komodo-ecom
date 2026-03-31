# komodo-order-reservations-api

Time-slot booking service for technician appointments and delivery windows.

## Port
`7063` (Orders block — see monorepo port allocation)

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Health check |
| `GET` | `/slots` | None | List available slots (`?date=YYYY-MM-DD&zone=<zone>`) |
| `GET` | `/slots/{date}` | None | Available slots for a specific date |
| `POST` | `/bookings` | JWT | Create/reserve a booking |
| `GET` | `/bookings/{id}` | JWT | Get booking by ID |
| `PUT` | `/bookings/{id}/cancel` | JWT | Cancel a booking |
| `PUT` | `/bookings/{id}/confirm` | JWT (internal) | Confirm a HELD booking after payment |

> **TODO:** Add `POST /internal/slots/sync` — internal endpoint for the external office scheduling
> system to push schedule changes into the DynamoDB read model.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `APP_NAME` | `komodo-order-reservations-api` |
| `PORT` | Service port (default `7063`) |
| `ENV` | `local` / `dev` / `staging` / `prod` |
| `LOG_LEVEL` | `debug` / `info` / `error` |
| `AWS_REGION` | AWS region |
| `AWS_ENDPOINT` | LocalStack endpoint for local dev |
| `AWS_SECRET_PREFIX` | Secrets Manager prefix |
| `AWS_SECRET_BATCH` | Batch secrets key |

### Secrets (via Secrets Manager)
| Key | Description |
|-----|-------------|
| `ORDER_RESERVATIONS_API_CLIENT_ID` | Service client ID |
| `ORDER_RESERVATIONS_API_CLIENT_SECRET` | Service client secret |
| `DYNAMODB_SLOTS_TABLE` | DynamoDB table name for slots read model |
| `DYNAMODB_BOOKINGS_TABLE` | DynamoDB table name for bookings |
| `RATE_LIMIT_RPS` | Rate limit requests/sec |
| `RATE_LIMIT_BURST` | Rate limit burst size |
| `BOOKING_HOLD_TTL_SECONDS` | **TODO (Option A):** TTL for HELD bookings during checkout |

## Run Commands

```bash
# Local (standalone)
make bootstrap

# Stop
make stop

# Lint
make lint

# Tests
make test
```

## Data Model

### Slots Table (`DYNAMODB_SLOTS_TABLE`)
Read model populated by the external office scheduling system (push-based sync).

| Key | Type | Description |
|-----|------|-------------|
| PK | `TechnicianID` | Partition key |
| SK | `SlotDateTime` | Sort key (ISO8601) |
| GSI | `SlotsByDate` | PK: `Date` (YYYY-MM-DD), SK: `SlotDateTime` |

**TODO:** Finalize table design — see `internal/repository/slot.go` for full notes.

### Bookings Table (`DYNAMODB_BOOKINGS_TABLE`)

| Key | Type | Description |
|-----|------|-------------|
| PK | `BookingID` | UUID |
| GSI | `BookingsByCustomer` | PK: `CustomerID`, SK: `CreatedAt` |
| GSI | `BookingsBySlot` | PK: `SlotDateTime+TechnicianID`, SK: `Status` |

**TODO:** Finalize table design — see `internal/repository/booking.go` for full notes.

## Open TODOs

- [ ] Decide checkout flow (Option A: pre-payment hold vs Option B: post-order scheduling) — see `internal/handlers/booking.go`
- [ ] Implement DynamoDB client + ConditionalWrite for double-booking prevention
- [ ] Add `POST /internal/slots/sync` for schedule push from external system
- [ ] Wire `PUT /bookings/{id}/confirm` to `order-api` (internal auth)
- [ ] Emit booking status-change events via `event-bus-api`
- [ ] Add `RangeOrderReservation = 45` to `komodo-forge-sdk-go/http/errors/ranges.go`
- [ ] Add `GET /me/bookings` for customer booking history (requires BookingsByCustomer GSI)
