# LocalStack Setup Guide

This monorepo uses a **centralized LocalStack** instance for local AWS service emulation across all APIs.

## Architecture

```
komodo-apis/
├── localstack/                      # Shared LocalStack setup
│   ├── docker-compose.yml          # LocalStack service
│   ├── init/                       # Auto-run initialization scripts
│   │   ├── 01-init-secretsmanager.sh
│   │   ├── 02-init-s3.sh
│   │   ├── 03-init-dynamo.sh
│   │   └── 04-init-aurora.sh
│   └── README.md
├── komodo-auth-api/
│   └── build/
│       └── docker-compose.dev.yaml  # Uses shared LocalStack
└── komodo-user-api/
    └── build/
        └── docker-compose.dev.yaml  # Uses shared LocalStack
```

## Quick Start

### 1. Start LocalStack

```bash
cd localstack
docker-compose up -d
```

This creates:
- LocalStack container on port 4566
- Docker network `komodo-network`
- All AWS resources (secrets, tables, buckets)

### 2. Verify LocalStack is Running

```bash
# Check container status
docker-compose ps

# View initialization logs
docker-compose logs -f

# Test AWS services
docker exec komodo-localstack awslocal s3 ls
docker exec komodo-localstack awslocal dynamodb list-tables
```

### 3. Start Your API(s)

```bash
# Auth API
cd ../komodo-auth-api/build
docker-compose -f docker-compose.dev.yaml up -d

# User API
cd ../../komodo-user-api/build
docker-compose -f docker-compose.dev.yaml up -d
```

## What Gets Created

### Secrets Manager
- `komodo-auth-api/dev/all-secrets` - JWT keys, OAuth credentials, IP lists
- `komodo-user-api/dev/all-secrets` - API client credentials, IP lists

### S3 Buckets
- `komodo-user-uploads-dev` - User file uploads (CORS enabled)
- `komodo-static-assets-dev` - Public static assets
- `komodo-backups-dev` - Backup storage

### DynamoDB Tables
- `komodo-users-dev` - User accounts (email GSI)
- `komodo-user-profiles-dev` - User profile data
- `komodo-sessions-dev` - Session management (user_id GSI, streams enabled)
- `komodo-oauth-tokens-dev` - OAuth tokens (user_id GSI)

### ElastiCache (Redis)
- Available at `localstack:6379` (no password required for dev)

## Development Workflow

### Standard Workflow
```bash
# 1. Start LocalStack (once)
cd localstack && docker-compose up -d

# 2. Start APIs as needed
cd ../komodo-auth-api/build
docker-compose -f docker-compose.dev.yaml up -d

# 3. View logs
docker-compose logs -f auth-api

# 4. Stop API when done
docker-compose down
```

### Full Stack Development
```bash
# Start everything
cd localstack && docker-compose up -d
cd ../komodo-auth-api/build && docker-compose -f docker-compose.dev.yaml up -d
cd ../../komodo-user-api/build && docker-compose -f docker-compose.dev.yaml up -d

# Stop everything
docker-compose down  # in each directory
cd ../../localstack && docker-compose down
```

## Accessing LocalStack Services

### From Your APIs (Docker)
```yaml
environment:
  - AWS_ENDPOINT=http://localstack:4566
  - ELASTICACHE_ENDPOINT=localstack:6379
```

### From Your Host Machine
```bash
# Install awslocal CLI wrapper
pip install awscli-local

# Use awslocal commands
awslocal s3 ls
awslocal dynamodb scan --table-name komodo-users-dev
awslocal secretsmanager get-secret-value --secret-id komodo-auth-api/dev/all-secrets
```

### Direct AWS CLI (alternative)
```bash
aws --endpoint-url=http://localhost:4566 s3 ls
```

## Common Tasks

### View Secrets
```bash
awslocal secretsmanager list-secrets
awslocal secretsmanager get-secret-value --secret-id komodo-auth-api/dev/all-secrets
```

### Check DynamoDB Tables
```bash
awslocal dynamodb list-tables
awslocal dynamodb describe-table --table-name komodo-users-dev
awslocal dynamodb scan --table-name komodo-users-dev --max-items 5
```

### Manage S3 Buckets
```bash
awslocal s3 ls
awslocal s3 ls s3://komodo-user-uploads-dev
awslocal s3 cp test.txt s3://komodo-user-uploads-dev/
```

### Test Redis
```bash
docker exec -it komodo-localstack redis-cli -h localhost -p 6379
# Then: SET test "hello" / GET test
```

## Troubleshooting

### LocalStack won't start
```bash
# Check port conflicts
lsof -i :4566

# View detailed logs
docker-compose logs localstack

# Restart fresh
docker-compose down -v
docker-compose up -d
```

### Initialization scripts didn't run
```bash
# Check script permissions
ls -la init/

# Make executable
chmod +x init/*.sh

# Restart LocalStack
docker-compose restart
```

### API can't connect to LocalStack
```bash
# Verify network exists
docker network ls | grep komodo-network

# Verify API is on the network
docker inspect komodo-auth-api-dev | grep -A 5 Networks

# Recreate network if needed
docker network create komodo-network
```

### Reset everything
```bash
# Stop all services
cd komodo-auth-api/build && docker-compose down
cd ../../komodo-user-api/build && docker-compose down
cd ../../localstack && docker-compose down -v

# Restart LocalStack
docker-compose up -d
```

## Tips

- **LocalStack data is ephemeral** - All data is lost when the container stops
- **Scripts are idempotent** - Safe to restart LocalStack anytime
- **Network is shared** - All APIs communicate through `komodo-network`
- **Port 4566** - Single endpoint for all AWS services
- **No credentials needed** - LocalStack doesn't require AWS keys for dev

## LocalStack Pro Features

Some services require LocalStack Pro:
- RDS/Aurora (script will skip if not available)
- Lambda
- ECS
- Advanced IAM

For local development without Pro, use:
- Docker PostgreSQL/MySQL containers instead of RDS
- Direct function calls instead of Lambda
- Docker containers instead of ECS

## Next Steps

1. ✅ LocalStack is running
2. ✅ APIs are configured to use LocalStack
3. Start building and testing locally!

For more details, see `localstack/README.md`
