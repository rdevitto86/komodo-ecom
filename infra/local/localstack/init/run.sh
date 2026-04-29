#!/bin/bash
# Single entry point executed by LocalStack ready.d.
# Sources individual init scripts in dependency order.
# Individual scripts live in ./scripts/ — LocalStack does not recurse into subdirectories.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")/scripts" && pwd)"

echo "=== Komodo LocalStack Init ==="

"$SCRIPT_DIR/init-secretsmanager.sh"
"$SCRIPT_DIR/init-s3.sh"
"$SCRIPT_DIR/init-dynamodb.sh"
"$SCRIPT_DIR/init-sqs.sh"
"$SCRIPT_DIR/init-redis.sh"

echo "=== Komodo LocalStack Init Complete ==="
