# Komodo — Deployment

CloudFormation-based infrastructure and ECS Fargate deployment for all environments.

---

## Environments

| Env | Trigger | Approval |
|-----|---------|----------|
| **local** | `local_demo.sh` / `docker compose` | None |
| **dev** | Push to `main` | None (auto) |
| **stg** | Dev deploy succeeds | Manual gate in GitHub Actions |
| **prod** | Manual workflow dispatch | Required reviewer approval |

---

## Structure

```
deploy/
├── cfn/
│   ├── infra.yaml          # Shared infra: VPC, DynamoDB, ElastiCache, ECR, Secrets Manager
│   ├── services.yaml       # ECS cluster, ALB, Fargate task defs + services
│   └── parameters/
│       ├── dev.json        # Dev environment parameter overrides
│       ├── stg.json        # Staging parameter overrides
│       └── prod.json       # Production parameter overrides
├── scripts/
│   ├── deploy-infra.sh     # Deploy/update the infra stack
│   └── deploy-services.sh  # Deploy/update the services stack
└── README.md (this file)
```

GitHub Actions workflows live at `../.github/workflows/`:

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `ci.yml` | Push/PR to main | Tests + docker build checks per changed service |
| `deploy-dev.yml` | CI passes on main | Auto-deploy changed services to DEV |
| `deploy-stg.yml` | Manual dispatch | Promote selected services to STG |
| `deploy-prod.yml` | Manual dispatch + approval | Promote selected services to PROD |
| `_deploy-service.yml` | Reusable (called by above) | Build → ECR push → ECS task def → service update |

### Per-Service Toggle

DEV: automatic — only services whose paths changed in the push are deployed.
STG/PROD: explicit checkboxes in the `workflow_dispatch` UI. All default to `false`.

### Required GitHub Secrets

| Secret | Description |
|--------|-------------|
| `AWS_ROLE_ARN_DEV` | OIDC role ARN for DEV account |
| `AWS_ROLE_ARN_STG` | OIDC role ARN for STG account |
| `AWS_ROLE_ARN_PROD` | OIDC role ARN for PROD account |
| `AWS_ACCOUNT_ID` | AWS account ID (used in ECR URI construction) |

PROD additionally requires the `production` GitHub Environment with required reviewers configured:
`Settings → Environments → production → Required reviewers`.

---

## Compute Strategy

All Komodo services run on **ECS Fargate** (not Lambda). The dual-port pattern
(public:anchor + internal:anchor+1) requires a persistent, multi-port runtime that
Lambda cannot provide.

| Exception | Compute | Reason |
|-----------|---------|--------|
| `komodo-analytics-collector-api` | Lambda | Write-only event sink, bursty, no dual-port |
| `komodo-core-features-api` | Lambda | Read-only flag evaluation, cacheable at edge |
| `komodo-core-entitlements-api` | Lambda | Read-only entitlement checks |

Cost controls on Fargate: Fargate Spot on dev/stg (up to 70% savings). Set desired
count to 0 on dev when idle. Single NAT Gateway (not per-AZ) on non-prod.

---

## Local Schema Reference

DynamoDB table definitions and ElastiCache config for local dev live in
`apis/localstack/init/`. The CloudFormation templates (`cfn/infra.yaml`)
are written to match those definitions exactly — keep them in sync when
adding or modifying tables.

---

## Prerequisites

- AWS CLI configured with appropriate IAM permissions
- ECR repositories created (handled by `infra.yaml` first run)
- Docker images built and pushed to ECR before `services.yaml` deploy

---

## Deploy Order

Always deploy infra before services on first run or when infra changes:

```bash
# 1. Deploy shared infrastructure
./scripts/deploy-infra.sh dev

# 2. Build + push Docker images to ECR (done by CI)

# 3. Deploy ECS services
./scripts/deploy-services.sh dev
```

---

## Stack Names

| Stack | Name pattern |
|-------|-------------|
| Infra | `komodo-infra-<env>` |
| Services | `komodo-services-<env>` |
