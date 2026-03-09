#!/bin/bash

echo "Initializing Redis..."

# Wait for Redis to be ready (with retries)
MAX_RETRIES=5
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" ping 2>/dev/null | grep -q "PONG"; then
    echo "✓ Redis is available at redis:6379"
    break
  fi
  
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "Waiting for Redis to be ready... (attempt $RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
  else
    echo "⚠️  Redis is not available after $MAX_RETRIES attempts"
    echo "   Redis container may not be running or not on the same network"
    exit 0
  fi
done

# Set some initial cache keys for testing
echo "Setting up test cache keys..."

redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" <<EOF
SET test:connection "Standalone Redis is working"
EXPIRE test:connection 3600
SET komodo:cache:initialized "true"
EXPIRE komodo:cache:initialized 86400
EOF

# Verify keys were set
if redis-cli -h redis -p 6379 -a "${AWS_ELASTICACHE_PASSWORD}" GET test:connection 2>/dev/null | grep -q "working"; then
  echo "✓ Test keys created successfully"
else
  echo "⚠️  Could not verify test keys"
fi

echo "Redis initialized successfully"
