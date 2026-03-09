# LocalStack Setup

Centralized LocalStack instance for local AWS emulation across all Komodo services.

## Quick Start

```bash
# From repo root — starts localstack + redis + auth-api
just up

# Verify LocalStack is healthy
docker logs localstack-main --tail 20

# Test AWS services
docker exec localstack-main awslocal s3 ls
docker exec localstack-main awslocal dynamodb list-tables
docker exec localstack-main awslocal secretsmanager list-secrets
```

Init scripts in `init/` run automatically on first start and seed all AWS resources.

---

## What Gets Created

### Secrets Manager
- `komodo/local/all-secrets` — JWT keys, OAuth credentials, DB endpoints, IP lists

### S3 Buckets
- `komodo-user-uploads-dev` — User file uploads (CORS enabled)
- `komodo-static-assets-dev` — Public static assets (SSR engine content)
- `komodo-backups-dev` — Backup storage

### DynamoDB Tables
- `komodo-users-dev` — User accounts (email GSI)
- `komodo-user-profiles-dev` — User profile data
- `komodo-sessions-dev` — Session management (user_id GSI, streams enabled)
- `komodo-oauth-tokens-dev` — OAuth tokens (user_id GSI)

### Redis
- Available at `komodo-redis:6379` within `komodo-network`
- Password: `test-password` (set in docker-compose)

---

## Accessing LocalStack from Host

```bash
# Install awslocal CLI wrapper
pip install awscli-local

# Examples
awslocal s3 ls
awslocal dynamodb list-tables
awslocal secretsmanager list-secrets
awslocal secretsmanager get-secret-value --secret-id komodo/local/all-secrets

# Or use standard AWS CLI with endpoint flag
aws --endpoint-url=http://localhost:4566 s3 ls
```

---

## Troubleshooting

### LocalStack won't start
```bash
lsof -i :4566          # check for port conflicts
just down-clean        # stop all + remove volumes, then just up again
```

### Init scripts didn't run
```bash
ls -la infra/local/localstack/init/   # check permissions
chmod +x infra/local/localstack/init/*.sh
just down && just up
```

### API can't connect
```bash
docker network ls | grep komodo-network    # verify network exists
docker inspect <container> | grep -A5 Networks
```

---

## Notes

- **Data is ephemeral** — lost when the container stops. Use `just down-clean` to fully reset.
- **Scripts are idempotent** — safe to restart LocalStack at any time.
- **Port 4566** — single endpoint for all AWS services.
- **No credentials needed** — LocalStack ignores AWS keys for local dev.
- **RDS/Aurora** requires LocalStack Pro — init script skips gracefully without it.
