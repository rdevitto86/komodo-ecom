# Komodo Infra — TODO

Priority guide: **[H]** = blocking local dev or AWS deployment · **[M]** = important, not blocking day-one · **[L]** = low priority (polish, monitoring, docs)

Sections are ordered by dependency — local dev must work before AWS deploy matters.

---

## Local Dev — Secrets Manager (LocalStack)
> All runnable APIs are now seeded. Remaining items are blocked on service scaffolding.

- [x] **[H]** Add `komodo-shop-items-api/local/all-secrets` to `01-init-secretsmanager.sh` — S3 bucket name, JWT keys, rate limits
- [x] **[H]** Add `komodo-cart-api/local/all-secrets` — ElastiCache endpoint/password/db, DynamoDB carts table, inventory-api URL, shop-items-api URL, JWT keys, hold TTL, guest TTL
- [x] **[H]** Add `komodo-shop-inventory-api/local/all-secrets` — DynamoDB inventory table, JWT keys, rate limits
- [x] **[H]** Add `komodo-order-api/local/all-secrets` — DynamoDB orders table, cart-api URL, payments-api URL, inventory-api URL, JWT keys
- [x] **[H]** Add `komodo-payments-api/local/all-secrets` — Stripe test keys, DynamoDB payments table, JWT keys
- [x] **[H]** Add `komodo-address-api/local/all-secrets` — address provider API key (stub value for local), JWT keys, rate limits
- [x] **[H]** Add `komodo-search-api/local/all-secrets` — Typesense host/API key (local Typesense or stub), JWT keys
- [x] **[M]** Add `komodo-support-api/local/all-secrets` — Anthropic API key, JWT keys, rate limits
- [ ] **[M]** Add `komodo-order-returns-api/local/all-secrets` — DynamoDB returns table, order-api URL, payments-api URL, JWT keys — **blocked: service not yet scaffolded**
- [x] **[M]** Add `komodo-order-reservations-api/local/all-secrets` — DynamoDB reservations table, JWT keys
- [x] **[H]** Add `komodo-event-bus-api/local/all-secrets` to `01-init-secretsmanager.sh` — JWT keys, `DYNAMO_EVENTS_TABLE=komodo-events`, `DYNAMO_SUBSCRIPTIONS_TABLE=komodo-event-subscriptions`, `DYNAMODB_ENDPOINT=http://host.docker.internal:4566`, `EVENT_TRANSPORT=dynamo`
- [x] **[M]** Add `komodo-communications-api/local/all-secrets` — email/SMS provider keys (stub for local), JWT keys
- [x] **[M]** Add `komodo-loyalty-api/local/all-secrets` — DynamoDB loyalty table, JWT keys
- [ ] **[M]** Add `komodo-reviews-api/local/all-secrets` — DynamoDB reviews table, JWT keys — **blocked: service not yet scaffolded**
- [ ] **[L]** Add `komodo-entitlements-api/local/all-secrets` and `komodo-features-api/local/all-secrets` once those services are scaffolded

---

## Local Dev — DynamoDB Tables (LocalStack)
> Only `komodo-users`, `komodo-sessions`, and `komodo-oauth-tokens` are created. All other service tables are TODO comments in `03-init-dynamodb.sh`.

- [x] **[H]** Add `komodo-events` table to `03-init-dynamodb.sh` — PK: `event_id` (S), SK: `domain` (S); enable Streams NEW_AND_OLD_IMAGES; TTL attribute `expires_at`
- [x] **[H]** Add `komodo-event-subscriptions` table to `03-init-dynamodb.sh` — PK: `event_type` (S), SK: `subscriber_url` (S), attr: `service_name`, `active`; seed with placeholder subscriber records
- [ ] **[H]** Add `komodo-carts` table to `03-init-dynamodb.sh` — PK: `CART#<userId>`, SK: `METADATA` | `ITEM#<itemId>`; no TTL; streams optional
- [ ] **[H]** Add `komodo-inventory` table — schema TBD in `apis/komodo-shop-inventory-api/docs/data-model.md` first; enable streams for CDC
- [ ] **[H]** Add `komodo-orders` table — schema TBD in `apis/komodo-order-api/docs/data-model.md` first; enable streams for CDC
- [ ] **[H]** Add `komodo-payments` table — schema TBD in `apis/komodo-payments-api/docs/data-model.md` first; enable streams
- [ ] **[M]** Add `komodo-returns` table — schema TBD in `apis/komodo-order-returns-api/docs/data-model.md`
- [ ] **[M]** Add `komodo-reservations` table — schema TBD in `apis/komodo-order-reservations-api/docs/data-model.md`
- [ ] **[M]** Add `komodo-support-sessions` table — replace in-memory storage in `support-api`
- [ ] **[L]** Add `komodo-loyalty`, `komodo-reviews` tables once those services are being implemented

---

## Local Dev — Docker Compose
> All runnable services have compose blocks. Remaining items are either scaffolding-blocked or optional.

- [x] **[H]** Add `address-api` service to `infra/local/docker-compose.yml` (port 7031)
- [x] **[H]** Add `cart-api` to compose (port 7043)
- [x] **[H]** Add `shop-inventory-api` to compose (port 7044)
- [x] **[M]** Add `support-api` to compose (port 7101)
- [x] **[M]** Add `order-reservations-api` to compose (port 7063)
- [x] **[M]** Add `event-bus-api` to compose (port 7002)
- [ ] **[L]** Add `entitlements-api` (7021), `features-api` (7022) compose blocks once those services are scaffolded
- [ ] **[L]** Add Typesense container to local compose so `search-api` can run fully locally

---

## CloudFormation — `infra.yaml` (DynamoDB + ECR + S3)
> Only 3 DynamoDB tables and 7 ECR repos are defined. Missing tables for 9 services and ECR repos for 8 services.

- [ ] **[H]** Add `komodo-carts` DynamoDB table resource — PAY_PER_REQUEST, streams NEW_AND_OLD_IMAGES, export StreamArn in Outputs
- [ ] **[H]** Add `komodo-inventory` DynamoDB table resource — streams enabled, export StreamArn
- [ ] **[H]** Add `komodo-orders` DynamoDB table resource — streams enabled, PITR enabled, export StreamArn
- [ ] **[H]** Add `komodo-payments` DynamoDB table resource — streams enabled, PITR enabled, export StreamArn
- [ ] **[M]** Add `komodo-returns` DynamoDB table resource
- [ ] **[M]** Add `komodo-reservations` DynamoDB table resource
- [ ] **[M]** Add `komodo-support-sessions` DynamoDB table resource
- [ ] **[H]** Add `komodo-events` DynamoDB table — PAY_PER_REQUEST, Streams NEW_AND_OLD_IMAGES, TTL on `expires_at`, export StreamArn
- [ ] **[H]** Add `komodo-event-subscriptions` DynamoDB table — PAY_PER_REQUEST, no streams
- [ ] **[H]** Wire `EventsTableStreamArn` in `event-pipeline.yaml` once table is in `infra.yaml`
- [ ] **[L]** Add `komodo-loyalty`, `komodo-reviews` DynamoDB table resources
- [ ] **[H]** Add ECR repositories for missing services: `komodo-cart-api`, `komodo-shop-inventory-api`, `komodo-order-api`, `komodo-payments-api`, `komodo-support-api`, `komodo-search-api`, `komodo-communications-api`, `komodo-order-returns-api`
- [ ] **[M]** Add ECR repos for: `komodo-loyalty-api`, `komodo-reviews-api`, `komodo-order-reservations-api`, `komodo-features-api`, `komodo-entitlements-api`
- [ ] **[M]** Enable streams + export StreamArn on `komodo-oauth-tokens` table (currently streams disabled; CDC needs it)
- [ ] **[M]** Add S3 bucket for email templates — versioning enabled, read-only bucket policy for `communications-api` task IAM role

---

## CloudFormation — `services.yaml` (ECS Task Definitions + Services)
> Only 4 of 15 services have task definitions, target groups, and ECS services defined. 11 are completely missing.

- [ ] **[H]** Add task definition + ECS service for `cart-api` (port 7043) — public-facing, needs ElastiCache + DynamoDB env vars
- [ ] **[H]** Add task definition + ECS service for `shop-inventory-api` (port 7044) — public-facing, DynamoDB env vars
- [ ] **[H]** Add task definition + ECS service for `order-api` (port 7061) — public-facing, DynamoDB + inter-service URLs
- [ ] **[H]** Add task definition + ECS service for `payments-api` (port 7071) — public-facing, DynamoDB + Stripe env vars
- [ ] **[H]** Add ALB listener rules for all 4 new public services above
- [ ] **[M]** Add task definition + ECS service for `address-api` (port 7031) — Lambda target (not ECS); add Lambda function resource or route to Lambda invocation
- [ ] **[M]** Add task definition + ECS service for `support-api` (port 7101)
- [ ] **[M]** Add task definition + ECS service for `search-api` (port 7042)
- [ ] **[M]** Add task definition + ECS service for `event-bus-api` (port 7002) — internal only, no ALB rule needed
- [ ] **[M]** Add task definition + ECS service for `order-returns-api` (port 7062)
- [ ] **[M]** Add task definition + ECS service for `order-reservations-api` (port 7063)
- [ ] **[M]** Add task definition + ECS service for `communications-api` (port 7081) — internal only
- [ ] **[L]** Add task definitions for `loyalty-api` (7091), `reviews-api` (7092) once implemented
- [ ] **[L]** Add task definitions for `entitlements-api` (7021), `features-api` (7022) once implemented
- [ ] **[M]** Add HTTPS listener + ACM certificate ARN parameter to ALB — currently only HTTP→HTTPS redirect is coded; HTTPS listener itself is a TODO
- [ ] **[M]** Add `address-api` and `payments-api` as Lambda functions (not ECS) per deployment strategy in root CLAUDE.md

---

## Deployment Strategy — Lambda vs Fargate
> Current architecture assumes Fargate/ECS with dual ports per runtime. Lambda deployment requires different approach.

- [ ] **[M]** Add deployment flag parameter (`DEPLOYMENT_TARGET: lambda|fargate`) to CloudFormation templates and deployment scripts
- [ ] **[M]** When `DEPLOYMENT_TARGET=lambda`, extract each binary into its own Lambda function instead of single multi-port ECS task
- [ ] **[M]** For services with both `/public` and `/private` endpoints (e.g., dual-port runtimes), create separate Lambda functions for public and private access since Lambda doesn't support multiple ports
- [ ] **[M]** Update CI/CD workflows to build and deploy Lambda functions when flag is set to lambda (requires separate packaging per function)
- [ ] **[M]** Add Lambda-specific resources to CloudFormation: IAM roles, API Gateway or ALB integration, function definitions per binary
- [ ] **[L]** Document cost and performance trade-offs between Lambda and Fargate deployment models

---

## CloudFormation — `event-pipeline.yaml` (CDC + SNS/SQS)
> SNS topics, SQS queues, and CDC Lambda are defined but gated on table stream ARNs that don't exist yet.

- [ ] **[H]** Wire `OrdersTableStreamArn` parameter — once `komodo-orders` table is added to `infra.yaml`, import its StreamArn here
- [ ] **[H]** Wire `PaymentsTableStreamArn` parameter — same pattern
- [ ] **[H]** Wire `CartTableStreamArn` parameter — same pattern
- [ ] **[H]** Wire `InventoryTableStreamArn` parameter — same pattern
- [ ] **[H]** Create `Dockerfile.cdc` in `apis/komodo-event-bus-api/` targeting `cmd/cdc` — referenced by the Lambda function resource but doesn't exist
- [ ] **[M]** Add EventBridge event bus + rules to `event-pipeline.yaml` — route CDC domain events to per-service SQS queues via EventBridge rules for flexible, per-domain fanout
- [ ] **[M]** Add SNS filter policies to subscriptions — currently all consumers receive all event types; scope each subscription to relevant event types only
- [ ] **[M]** Add `komodo-support-events` SNS topic + consumer queues (escalation → communications)
- [ ] **[L]** Narrow CDC Lambda IAM role DynamoDB stream ARN from wildcard (`*`) to explicit table ARNs once all tables are defined

---

## CI/CD — GitHub Actions
> All workflows are disabled (manual `workflow_dispatch` only). Docker build and deploy matrices are incomplete.

- [ ] **[H]** Fix hard-coded `AWS_REGION: us-east-1` in `_deploy-service.yml` → `us-east-2`
- [ ] **[H]** Add missing services to CI Docker build matrix in `ci.yml`: `cart-api`, `shop-inventory-api`, `order-api`, `payments-api`, `support-api`, `search-api`
- [ ] **[H]** Add missing services to deploy matrix in `deploy-dev.yml` and `deploy-prod.yml`: `address-api`, `cart-api`, `shop-inventory-api`, `order-api`, `payments-api`, `support-api`, `order-returns-api`, `order-reservations-api`, `event-bus-api`
- [ ] **[M]** Configure GitHub OIDC: create AWS IAM role with trust policy for the GitHub repo; add `AWS_ACCOUNT_ID` and `AWS_ROLE_ARN_DEV`/`AWS_ROLE_ARN_PROD` as GitHub secrets
- [ ] **[M]** Re-enable CI auto-trigger — uncomment `on.pull_request` and `on.push` blocks in `ci.yml`
- [ ] **[M]** Re-enable deploy-dev auto-trigger — uncomment `on.workflow_run` in `deploy-dev.yml` once CI is stable
- [ ] **[M]** Add CDC Lambda image build + push step to CI (needs `Dockerfile.cdc` first)
- [ ] **[M]** Add `deploy-event-pipeline.sh` script and a corresponding CI step to deploy `event-pipeline.yaml` after `infra.yaml`
- [ ] **[L]** Add post-deploy smoke test step (hit `/health` on each newly deployed service)
- [ ] **[L]** Add rollback step on ECS deployment failure
- [ ] **[L]** Add `golangci-lint` step to CI

---

## Production Secrets Seeding
> No automation exists to seed AWS Secrets Manager for real environments. ECS tasks will fail on first deploy without this.

- [ ] **[H]** Create `infra/deploy/scripts/seed-secrets.sh` — idempotent script that creates/updates all required secrets in AWS Secrets Manager for a given environment; reads from a local `secrets.env` file (gitignored)
- [ ] **[H]** Document required secret keys per service in `infra/deploy/secrets-manifest.md` — single reference for ops to know what to populate before deploying
- [ ] **[M]** Add secret ARN outputs to `infra.yaml` so `services.yaml` task definitions can reference them by ARN instead of path
- [ ] **[L]** Define secret rotation policy for JWT keys and database credentials

---

## EC2 Deploy
> EC2 path covers only 3 services. No nginx config exists in the repo.

- [ ] **[M]** Add `nginx.conf` to `infra/deploy/ec2/` — reverse proxy config for TLS termination and routing to service ports
- [ ] **[M]** Add remaining core services to `infra/deploy/ec2/docker-compose.yaml`: `cart-api`, `shop-inventory-api`, `order-api`, `payments-api`
- [ ] **[M]** Document Certbot SSL setup steps in `ec2/` — setup.sh installs certbot but no guidance on running it
- [ ] **[L]** Add `komodo-support-api`, `komodo-search-api` to EC2 compose

---

## Networking & Security
> VPC and security groups are solid. These are hardening items.

- [ ] **[M]** Add VPC endpoints for Secrets Manager, DynamoDB, and S3 in `infra.yaml` — reduces NAT Gateway traffic and keeps AWS API calls off the public internet
- [ ] **[M]** Enable PITR (Point-in-Time Recovery) on all DynamoDB tables with financial data (`komodo-orders`, `komodo-payments`) — currently only enabled on `komodo-users`
- [ ] **[L]** Add Network ACLs with explicit deny rules as a second layer below security groups
- [ ] **[L]** Add bastion host or SSM Session Manager access for EC2 instances — no documented way to SSH in currently

---

## Data Lake & Analytics Pipeline
> Phase 1: Events API writes directly to S3 on publish. Phase 2: Add Kinesis as the pipeline backbone when volume warrants — gate behind cost evaluation. Athena handles batch reporting over S3 data.

- [ ] **[M]** Provision S3 data lake bucket in `infra.yaml` — versioning enabled, lifecycle policy (e.g. 90-day transition to Glacier), bucket policy scoped to Events API task IAM role; this is the archive destination for all published events
- [ ] **[M]** Wire Events API to write event payloads to S3 on publish — partition by domain and date (`s3://komodo-data-lake/events/{domain}/{YYYY}/{MM}/{DD}/`); Phase 1 path bypasses Kinesis to avoid cost
- [ ] **[M]** Set up Athena workgroup + database over the data lake bucket — enables monthly/quarterly ad-hoc reporting and analytics queries against the full event history; output results to a separate S3 results bucket
- [ ] **[L]** Phase 2 — Add Kinesis Data Firehose between Events API and S3 when event volume warrants — Firehose buffers, batches, and delivers to S3; Events API publishes to Kinesis stream instead of writing S3 directly; gate this behind a cost/volume threshold review
- [ ] **[L]** Phase 2 — Evaluate Kinesis Data Streams for real-time consumers (Statistics API, external subscribers) — if bi-directional streaming is needed beyond SNS/SQS fanout, Kinesis Streams adds per-shard consumer support; document tradeoffs before implementing

---

## Monitoring & Observability
> Only DLQ depth alarms exist. No service health, error rate, or latency monitoring.

- [ ] **[M]** Add CloudWatch alarms for ECS service CPU and memory utilization (alert at 80%)
- [ ] **[M]** Add ALB 5xx error rate alarm per target group
- [ ] **[M]** Add ALB target response time alarm (P99 > 2s)
- [x] ~~Add telemetry ingestion endpoint (Lambda or API Gateway)~~ — V1 strategy: all services write structured JSON logs to their own CloudWatch stream via `slog`. UI client logs route through the SvelteKit BFF `/api/log` route (no separate Lambda needed). Data lake is a future concern.
- [ ] **[L]** Create a CloudWatch dashboard for the core purchase flow: auth → cart → checkout → order
- [ ] **[L]** Add CloudWatch Synthetics canary for `/health` checks on all public services
- [ ] **[L]** Set up CloudWatch log metric filters to alert on `"level":"error"` in structured logs

---

## Testing
- [ ] **[L]** Write LocalStack init integration test — script that starts LocalStack, runs all init scripts, and asserts all tables/secrets/buckets were created correctly
- [ ] **[L]** Write CloudFormation template linting step (`cfn-lint`) in CI for `infra.yaml`, `services.yaml`, `event-pipeline.yaml`
- [ ] **[L]** Add `checkov` or `cfn-guard` security scan on CloudFormation templates to catch IAM wildcards, unencrypted resources, etc.

---

## Documentation & Diagrams
- [ ] **[L]** Create architecture and data-flow diagrams in `docs/` — service topology, checkout flow, CDC event pipeline, auth flow
