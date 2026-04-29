# Komodo Platform — Route Audit

**Auth model (V1 — AWS)**
- External traffic: API Gateway (HTTP API) → Lambda / Fargate. Gateway JWT authorizer validates user tokens.
- Service-to-service: VPC-direct HTTP with M2M JWTs (client credentials via auth-api). Services validate token type and branch behavior accordingly.
- JWKS endpoint is env config (`JWT_JWKS_URL`) — swappable without code changes when issuer migrates.

**Annotations**
- `← auth` — user JWT required
- `← service` — service JWT required (M2M, client credentials); VPC-routed
- `← guest` — JWT optional; fallback via query param or session cookie
- `← ext` — external third-party webhook (Stripe, carrier, etc.); arrives via public internet, signature-validated not JWT-validated
- `[done]` implemented · `[stub]` registered, returns 501 · `[spec-only]` specced, no code · `[TODO]` in TODO.md, not yet specced
- `[?]` open question · `[V2]` not blocking ecom MVP
- Health endpoints omitted — every service has `GET /health`

---

## Identity & Security

### komodo-auth-api · Fargate · 7011 / 7012 · `[done]`

```
# public (7011)
GET  /.well-known/jwks.json       ← pub
POST /oauth/token                 ← pub
GET  /oauth/authorize             ← pub
POST /oauth/revoke                ← pub
POST /otp/request                 ← pub
POST /otp/verify                  ← pub

# private (7012)
POST /oauth/introspect            ← priv
POST /token/validate              ← priv
GET  /clients                     ← priv
GET  /clients/{id}                ← priv
```

---

## Core Platform

### komodo-platforms-api · Lambda · 7022 · `[spec-only]`

```
# private only
GET /access?appId={appId}                             ← priv
PUT /access # [TODO] add body with role, toggles, entitlements  ← priv
```

### komodo-ai-guardrails-api · Lambda · 7023 · `[spec-only]` (Python FastAPI, not yet created)

```
# private only
POST /moderate                    ← priv
```

---

## Address & Geo

### komodo-address-api · Lambda · 7031 · `[stub]` (provider calls not wired)

```
# public only
POST /addresses/validate          ← pub
POST /addresses/normalize         ← pub
POST /addresses/geocode           ← pub
```

---

## Commerce & Catalog

### komodo-shop-items-api · Fargate · 7041 · `[done]`

```
# public only
GET  /item/inventory?sort={sort}&filters={filters}    ← pub
GET  /item/{sku}                  ← pub
GET  /services/repair?sort={sort}&filters={filters}   ← pub
GET  /services/repair/{id}        ← pub
```

### komodo-search-api · Fargate · 7042 · `[stub]` (Typesense not wired)

```
# public
GET /search                       ← pub

# private
POST   /index/sync                ← priv
DELETE /index                     ← priv
```

### komodo-cart-api · Fargate · 7043 · `[done]`

```
# public only
GET    /cart                      ← pub
DELETE /cart                      ← pub
POST   /cart/merge                ← pub
POST   /cart/items                ← pub
PUT    /cart/items/{itemId}       ← pub
DELETE /cart/items/{itemId}       ← pub
POST   /cart/checkout             ← pub
GET    /cart/saved                ← pub
POST   /cart/saved                ← pub
DELETE /cart/saved/{itemId}       ← pub
POST   /cart/saved/{itemId}/move  ← pub
```

### komodo-shop-inventory-api · Fargate · 7044 · `[stub]` (DynamoDB ops are todo!())

```
# private only
GET    /stock                     ← priv
GET    /stock/{sku}               ← priv
POST   /stock/{sku}/reserve       ← priv
DELETE /stock/{sku}/reserve/{holdId} ← priv
POST   /stock/{sku}/confirm       ← priv
POST   /stock/{sku}/restock       ← priv
```

### komodo-shop-promotions-api · Lambda · 7045 · `[spec-only]`

```
# public
POST /validate                    ← pub
GET  /promotions                  ← pub
GET  /promotions/{promoId}        ← pub

# private
POST   /promotions                ← priv
PUT    /promotions/{promoId}      ← priv
DELETE /promotions/{promoId}      ← priv
```

---

## User & Profile

### komodo-user-api · Fargate · 7051 / 7052 · `[done]`

```
# public (7051)
GET    /profile                   ← pub
POST   /profile                   ← pub
PUT    /profile                   ← pub
DELETE /profile                   ← pub
GET    /addresses                 ← pub
POST   /addresses                 ← pub
PUT    /addresses/{id}            ← pub
DELETE /addresses/{id}            ← pub
GET    /payments                  ← pub
PUT    /payments                  ← pub
DELETE /payments/{id}             ← pub
GET    /preferences               ← pub
PUT    /preferences               ← pub
DELETE /preferences               ← pub
GET    /wishlist                  ← pub
POST   /wishlist/items            ← pub
DELETE /wishlist/items/{itemId}   ← pub
GET    /wishlist/availability     ← pub
POST   /wishlist/move-to-cart     ← pub

# private (7052)
GET /users/{id}                   ← priv
GET /users/{id}/addresses         ← priv
GET /users/{id}/preferences       ← priv
GET /users/{id}/payments          ← priv
```

---

## Orders & Fulfillment

### komodo-order-api · Fargate · 7061 / 7062 · `[done]` (returns stubbed)

```
# public (7061)
POST   /orders                    ← pub
GET    /orders                    ← pub
GET    /orders/{orderId}          ← pub
POST   /orders/{orderId}/cancel   ← pub
GET    /orders/returns            ← pub
POST   /orders/returns            ← pub
GET    /orders/returns/{returnId} ← pub
DELETE /orders/returns/{returnId} ← pub

# private (7062)
PUT /returns/{returnId}/approve   ← priv
PUT /returns/{returnId}/receive   ← priv
PUT /returns/{returnId}/reject    ← priv
```

### komodo-order-reservations-api · Fargate · 7063 · `[stub]` (repo is all stubs)

```
# public
GET  /slots?date={date}&serviceType={serviceType}&locationId={locationId}  ← pub
GET  /reservations                ← pub
POST /reservations                ← pub
GET  /reservations/{id}           ← pub
PUT  /reservations/{id}           ← pub
PUT  /reservations/{id}/confirm   ← pub

# private
POST /slots/sync                  ← priv
```

### komodo-shipping-api · Lambda · 7064 · `[TODO]` (service not yet created)

```
# public
GET /shipments/{shipmentId}       ← pub
POST /labels/outbound             ← pub
POST /labels/inbound              ← pub

# external
POST /webhooks/carrier            ← ext
```

---

## Payments

### komodo-payments-api · Lambda · 7071 · `[stub]` (DynamoDB + Stripe are todo!())

```
# public
GET    /payments                           ← pub
GET    /payments/methods                   ← pub
POST   /payments/methods                   ← pub
DELETE /payments/methods/{methodId}        ← pub
GET    /payments/plans?status={status}     ← pub
GET    /payments/plans/{planId}            ← pub
POST   /payments/charge                    ← pub
POST   /payments/refund                    ← pub
GET    /payments/{chargeId}                ← pub
POST   /payments/plans/{planId}/execute    ← pub

# private
POST   /payments/plans                     ← priv
DELETE /payments/plans/{planId}            ← priv

# external
POST /webhooks/payments                    ← ext
```

---

## Communications

### komodo-communications-api · Lambda · 7081 · `[spec-only]` (main.go is empty)

```
# private only
POST /send/email                  ← priv
POST /send/sms                    ← priv
POST /send/push                   ← priv
```

---

## Loyalty & Social

### komodo-loyalty-api · Lambda · 7091 · `[stub]` (all handlers return 501)

```
# public only
GET    /loyalty                   ← pub
POST   /loyalty/redeem            ← pub
POST   /reviews                   ← pub
PUT    /reviews/{reviewId}        ← pub
DELETE /reviews/{reviewId}        ← pub
GET    /items/{itemId}/reviews    ← pub
```

---

## Support & CX

### komodo-support-api · Fargate · 7101 · `[done]` (in-memory — needs DynamoDB before prod)

```
# public only
POST   /chat/session              ← pub
GET    /chat/session              ← pub
POST   /chat/message              ← pub
GET    /chat/history?session={sessionId} # [TODO] wildcard for fetch all messages for a user  ← pub
DELETE /chat/history              ← pub
POST   /chat/escalate             ← pub
```

---

## Analytics & Discovery

### komodo-statistics-api · Fargate · 7111 / 7112 · `[stub]` (SQLite placeholder, no events consumed)

```
# public (7111)
GET  /stats/items/trending                 ← pub
GET  /stats/items/{itemId}/in-cart         ← pub
GET  /stats/items/{itemId}/recently-bought ← pub
GET  /stats/items/{itemId}/frequently-bought-with  ← pub
POST /events                               ← pub

# private (7112)
GET /stats/dashboard?from={fromDate}&to={toDate}&granularity={day|week|month}&category={category}  ← priv
```

### komodo-insights-api · Lambda · 7114 · `[stub]` (LLM provider not initialized)

```
# public only
GET /items/trending               ← pub
GET /items/{itemId}/summary       ← pub
GET /items/{itemId}/sentiment     ← pub
```

---

## Infrastructure

### komodo-events-api · Fargate · 7002 · `[done]` (local; SNS/SQS not deployed)

```
# private only
POST /publish                     ← priv
POST /subscribe                   ← priv
```

### komodo-ssr-engine-svelte · Fargate · 7003 · `[spec-only]`

```
# private only
POST /render                      ← priv
POST /admin/content/invalidate    ← priv
POST /admin/content/upsert        ← priv
```
