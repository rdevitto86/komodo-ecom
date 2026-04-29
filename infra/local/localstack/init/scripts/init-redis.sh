#!/bin/bash
# Pings the Redis container (not LocalStack) and seeds test keys.
# redis-cli is available because the LocalStack entrypoint installs redis-tools.

echo "Initializing Redis..."

MAX_RETRIES=5
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" ping 2>/dev/null | grep -q "PONG"; then
    echo "Redis is available at redis:6379"
    break
  fi

  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "Waiting for Redis... (attempt $RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
  else
    echo "WARNING: Redis not available after $MAX_RETRIES attempts"
    exit 0
  fi
done

redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" <<EOF
SET test:connection "Standalone Redis is working"
EXPIRE test:connection 3600
SET komodo:cache:initialized "true"
EXPIRE komodo:cache:initialized 86400
EOF

if redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" GET test:connection 2>/dev/null | grep -q "working"; then
  echo "Test keys created successfully"
else
  echo "WARNING: Could not verify test keys"
fi

echo "Redis initialized successfully"
