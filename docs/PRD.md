# Komodo — Product Requirements Document

> Phased PRD covering the journey from a small personal project to a modern, multi-vertical ecommerce and operations platform. This document outlines the vision, goals, and roadmap for the Komodo platform. It is owned by the project owner and updated as the product matures.

---

## 1. Executive Summary

<!-- FILL IN: 3-5 sentences. What Komodo is, who it serves today, where it is heading. Frame the through-line from MVP storefront to multi-vertical platform. Keep it grounded — this is a real business, not a thought experiment. -->

The Komodo ecommerce monorepo is a scalable ecommerce platform built with a microservices architecture and modern web technologies.
It currently serves as a personal project to demonstrate design, planning, and engineering skills, with te potential to serve as a primary web-based shop for
3D-printed goods, machinied parts, electronic devices, and associated services.

---

## 2. Problem Statement

<!-- FILL IN: What problem does Komodo solve, and for whom? Why does the current market fail to solve it well enough? Why is the owner the right person to build it? Two short paragraphs. -->

---

## 3. Goals & Non-Goals

| Type | Statement |
|------|-----------|
| Goal | <!-- FILL IN: e.g. "Ship a working storefront for 3D-printed goods within X months at < $Y/mo run cost" --> |
| Goal | <!-- FILL IN --> |
| Goal | <!-- FILL IN --> |
| Non-Goal | <!-- FILL IN: e.g. "Building a generic Shopify competitor" --> |
| Non-Goal | <!-- FILL IN --> |

---

## 4. User Personas

| Persona | Description | Primary Needs |
|---------|-------------|---------------|
| <!-- e.g. Hobbyist Buyer --> | <!-- FILL IN --> | <!-- FILL IN --> |
| <!-- e.g. Custom Order Client --> | <!-- FILL IN --> | <!-- FILL IN --> |
| <!-- e.g. Future B2B Buyer (Ag/Warehouse) --> | <!-- FILL IN --> | <!-- FILL IN --> |
| Owner / Operator | Single operator running storefront, fulfillment, and support during MVP. | Low-friction admin, observable system, predictable cost. |

---

## 5. Phase 1 — MVP

<!-- FILL IN: 2-3 sentences framing what "done" looks like for the MVP. Anchor to the current state — frontend live (static), backend services scaffolded, EC2 docker-compose as the first deploy target. -->

### 5.1 Scope & Feature List

| Feature | Service(s) | Priority |
|---------|------------|----------|
| Browse catalog | shop-items-api, search-api, ui | P0 |
| Product detail page | shop-items-api, reviews-api, ui | P0 |
| Cart | cart-api, ui | P0 |
| Account / sign-in | auth-api, user-api, ui | P0 |
| Checkout & payment | order-reservations-api, payments-api, order-api, ui | P0 |
| Order confirmation email | communications-api, event-bus-api | P0 |
| Order history | order-api, ui | P1 |
| Address management | address-api, ui | P1 |
| Wishlist | user-wishlist-api, ui | P2 |
| Promo codes | shop-promotions-api | P2 |
| Customer support contact | support-api | P2 |
| Returns | order-returns-api | P2 |
| <!-- FILL IN: any other MVP features --> | | |

### 5.2 Success Metrics

<!-- FILL IN: pick 3-5 measurable outcomes. e.g. "first paid order shipped", "p95 page load < Xms", "monthly infra cost < $Y", "checkout conversion > Z%". Keep them honest — early-stage metrics, not vanity numbers. -->

### 5.3 Out of Scope (Phase 1)

- Multi-tenant / B2B accounts
- Loyalty program activation (service scaffolded, not active)
- SSR engine activation (service exists, dormant until perf demands it)
- Inventory forecasting and supplier integrations
- <!-- FILL IN: anything else explicitly deferred -->

---

## 6. Phase 2 — Growth

<!-- FILL IN: 2-3 sentences on what triggers the move from Phase 1 to Phase 2 (e.g. revenue threshold, order volume, hiring first contractor). -->

### 6.1 New Capabilities

- Loyalty program (loyalty-api activation)
- Reviews and social proof at scale (reviews-api hardening)
- Promotions engine (shop-promotions-api expansion)
- Insights and merchandising analytics (insights-api)
- Feature flagging for staged rollouts (features-api)
- Entitlements model for gated content / membership tiers (entitlements-api)
- <!-- FILL IN: any growth-phase features specific to Komodo's roadmap -->

### 6.2 Infrastructure Evolution

- **Compute migration:** EC2 docker-compose → ECS Fargate using existing CloudFormation templates in `infra/deploy/cfn/`. No code changes required — `cmd/public` and `cmd/private` entrypoints already align with Fargate task definitions.
- **CI/CD reactivation:** Re-enable GitHub Actions auto-triggers (currently `workflow_dispatch` only). Uncomment `on:` blocks in `ci.yml` and `deploy-dev.yml` once backend is live in dev.
- **Frontend adapter switch:** SvelteKit `adapter-static` → `adapter-node` once a real backend is wired. Commented import already present in `ui/svelte.config.js`.
- **SSR engine activation:** Bring `komodo-ssr-engine-svelte` (port 7003) online for SEO-critical and performance-sensitive pages.
- **Search backend:** Wire Typesense behind `search-api` (currently scaffold).
- **Inventory hardening:** Promote Rust `shop-inventory-api` from V2 experiment to primary.

### 6.3 Major Hurdles & Migrations

> **Callout — known hard transitions.** These are the changes most likely to bite. Each needs its own migration plan before being attempted.
>
> - **EC2 → Fargate cutover:** session affinity, secret rotation, log shipping, and health-check tuning all change shape.
> - **adapter-static → adapter-node:** routing, SSR data fetching, and CDN cache strategy all shift.
> - **Payments V1 (any) → Rust payments-api:** dual-write / shadow-traffic strategy required to avoid charge anomalies.
> - **Inventory V1 → Rust shop-inventory-api:** event-driven sync with search and catalog must be locked down before flip.
> - <!-- FILL IN: other migrations the owner anticipates -->

---

## 7. Phase 3 — Scale / B2B

<!-- FILL IN: 2-3 sentences on the long-range vision. Multi-vertical (Agriculture, Warehousing tech), B2B contracts, possibly white-label. Anchor it to why the platform's current architecture supports this — domain-bounded services, event-driven core, cleanly separated public/private entrypoints. -->

### 7.1 New Product Lines

| Line | Description | Net-new services / capabilities |
|------|-------------|---------------------------------|
| Agriculture tech | <!-- FILL IN: sensors, telemetry, equipment? --> | <!-- FILL IN --> |
| Warehousing tech | <!-- FILL IN: inventory automation, asset tracking? --> | <!-- FILL IN --> |

### 7.2 B2B Considerations

- **Account model:** Organization-level accounts with role-based access (procurement, approver, admin). Likely an extension of `user-api` + `entitlements-api`.
- **Pricing:** Contract pricing, net terms, quote-to-order workflows. New service or expansion of `shop-promotions-api` + `order-api`.
- **Compliance:** SOC 2 readiness, data residency, audit logging hardening (root `docs/` already houses audit-logging conventions).
- **Integrations:** EDI, ERP connectors, procurement portals (Coupa, Ariba, etc.).
- <!-- FILL IN: vertical-specific B2B needs -->

### 7.3 Infrastructure at Scale

> **Callout — scale architecture.** At this phase the platform is multi-region capable, with per-vertical isolation where regulation or SLA demands it.
>
> - Multi-account AWS organization, per-vertical workload accounts.
> - Read-replica / global table strategy for DynamoDB where appropriate.
> - Dedicated event backbone (EventBridge or MSK) replacing in-cluster `event-bus-api` for cross-vertical fan-out.
> - Observability stack hardened (centralized logs, metrics, traces — likely OpenTelemetry across all services).
> - <!-- FILL IN: other scale concerns -->

---

## 8. Risks & Open Questions

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Single-operator bandwidth becomes the bottleneck before revenue justifies a hire | High | <!-- FILL IN --> |
| EC2 → Fargate migration introduces regressions in auth or payment paths | Medium | Shadow traffic + canary rollout per service |
| Payments V2 (Rust) ships with a defect that affects real money | Medium | Dual-run against V1 in shadow mode; reconciliation job |
| DynamoDB access patterns ossify before product-market fit | Medium | Keep `data-model.md` per service current; revisit GSIs quarterly |
| Premature B2B work fragments the codebase before B2C is profitable | Medium | Hard gate Phase 3 work behind Phase 2 revenue milestone |
| <!-- FILL IN: risk --> | | |

**Open Questions**

- <!-- FILL IN: e.g. "Does Komodo own fulfillment or partner out?" -->
- <!-- FILL IN: e.g. "Self-hosted Typesense vs managed search?" -->
- <!-- FILL IN: e.g. "When does the SSR engine become net positive vs added complexity?" -->

---

## 9. Appendix / References

- `CLAUDE.md` — monorepo conventions, service registry, port allocation
- `docs/HLA.md` — high-level architecture
- `docs/LLA.md` — low-level architecture (container topology, runtime flows)
- `apis/<service>/docs/` — per-service architecture, design decisions, data models
- `infra/deploy/cfn/` — CloudFormation templates for ECS Fargate scale-up path
- <!-- FILL IN: external references, market research, competitor notes -->
