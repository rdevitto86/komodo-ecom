#!/bin/bash

# LocalStack initialization script for Aurora/RDS
# Note: RDS is a LocalStack Pro feature
# This script is a placeholder for future use

echo "Aurora/RDS initialization..."

# Check if LocalStack Pro is available
if awslocal rds describe-db-instances 2>&1 | grep -q "not available"; then
  echo "⚠️  RDS/Aurora is a LocalStack Pro feature"
  echo "   Skipping RDS initialization"
  echo "   For local development, consider using:"
  echo "   - Docker PostgreSQL container"
  echo "   - Docker MySQL container"
  exit 0
fi

# If Pro is available, create RDS instances
echo "Creating Aurora PostgreSQL cluster..."
awslocal rds create-db-cluster \
  --db-cluster-identifier komodo-aurora-dev \
  --engine aurora-postgresql \
  --engine-version 13.7 \
  --master-username komodo_admin \
  --master-user-password komodo_dev_password \
  --database-name komodo_dev \
  2>/dev/null || echo "Aurora cluster already exists"

awslocal rds create-db-instance \
  --db-instance-identifier komodo-aurora-instance-dev \
  --db-cluster-identifier komodo-aurora-dev \
  --engine aurora-postgresql \
  --db-instance-class db.t3.medium \
  2>/dev/null || echo "Aurora instance already exists"

echo "Aurora/RDS initialization complete"
