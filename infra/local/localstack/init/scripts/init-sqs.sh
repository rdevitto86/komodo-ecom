#!/bin/bash
# Creates the SNS topics and SQS queues that mirror the production event-pipeline.yaml topology.
# Naming follows the same convention as CFN but uses "local" instead of an environment suffix.
#
# Topology (sourced from infra/deploy/cfn/event-pipeline.yaml):
#   SNS FIFO topics  (1 per domain): order, user, payment, cart, inventory
#   SQS FIFO queues  (1 per subscriber): each with a paired DLQ
#   SQS standard     CDC Lambda failure queue (non-FIFO, used as OnFailure destination)
#   SNS→SQS subs     raw message delivery, no filter policies yet
#
# NOTE: none of this is active locally — event-bus-api uses EVENT_TRANSPORT=dynamo.
# These resources exist so future SQS consumers can be tested without infra changes.

echo "Initializing SQS/SNS in LocalStack..."

sleep 1

# ── SNS FIFO Topics ───────────────────────────────────────────────────────────

echo "Creating SNS FIFO topics..."

for DOMAIN in order user payment cart inventory; do
  awslocal sns create-topic \
    --name "komodo-${DOMAIN}-events-local.fifo" \
    --attributes FifoTopic=true,ContentBasedDeduplication=false \
    2>/dev/null || echo "Topic komodo-${DOMAIN}-events-local.fifo already exists"
done

# ── SQS FIFO DLQs ─────────────────────────────────────────────────────────────

echo "Creating SQS FIFO DLQs..."

DLQ_NAMES=(
  "komodo-order-events-inventory-local-dlq.fifo"
  "komodo-order-events-communications-local-dlq.fifo"
  "komodo-order-events-loyalty-local-dlq.fifo"
  "komodo-payment-events-order-local-dlq.fifo"
  "komodo-payment-events-communications-local-dlq.fifo"
  "komodo-user-events-loyalty-local-dlq.fifo"
  "komodo-user-events-communications-local-dlq.fifo"
)

for DLQ in "${DLQ_NAMES[@]}"; do
  awslocal sqs create-queue \
    --queue-name "$DLQ" \
    --attributes FifoQueue=true,MessageRetentionPeriod=1209600 \
    2>/dev/null || echo "Queue $DLQ already exists"
done

# ── SQS FIFO Subscriber Queues ────────────────────────────────────────────────

echo "Creating SQS FIFO subscriber queues..."

declare -A QUEUE_TO_DLQ=(
  ["komodo-order-events-inventory-local.fifo"]="komodo-order-events-inventory-local-dlq.fifo"
  ["komodo-order-events-communications-local.fifo"]="komodo-order-events-communications-local-dlq.fifo"
  ["komodo-order-events-loyalty-local.fifo"]="komodo-order-events-loyalty-local-dlq.fifo"
  ["komodo-payment-events-order-local.fifo"]="komodo-payment-events-order-local-dlq.fifo"
  ["komodo-payment-events-communications-local.fifo"]="komodo-payment-events-communications-local-dlq.fifo"
  ["komodo-user-events-loyalty-local.fifo"]="komodo-user-events-loyalty-local-dlq.fifo"
  ["komodo-user-events-communications-local.fifo"]="komodo-user-events-communications-local-dlq.fifo"
)

for QUEUE in "${!QUEUE_TO_DLQ[@]}"; do
  DLQ_NAME="${QUEUE_TO_DLQ[$QUEUE]}"
  DLQ_URL=$(awslocal sqs get-queue-url --queue-name "$DLQ_NAME" --query QueueUrl --output text 2>/dev/null || true)
  DLQ_ARN=$(awslocal sqs get-queue-attributes \
    --queue-url "$DLQ_URL" \
    --attribute-names QueueArn \
    --query Attributes.QueueArn --output text 2>/dev/null || true)

  awslocal sqs create-queue \
    --queue-name "$QUEUE" \
    --attributes "FifoQueue=true,ContentBasedDeduplication=false,VisibilityTimeout=300,RedrivePolicy={\"deadLetterTargetArn\":\"${DLQ_ARN}\",\"maxReceiveCount\":\"3\"}" \
    2>/dev/null || echo "Queue $QUEUE already exists"
done

# ── CDC Lambda Failure Queue (standard, non-FIFO) ─────────────────────────────

echo "Creating CDC Lambda failure queue..."

awslocal sqs create-queue \
  --queue-name "komodo-event-bus-cdc-failures-local" \
  --attributes MessageRetentionPeriod=1209600 \
  2>/dev/null || echo "Queue komodo-event-bus-cdc-failures-local already exists"

# ── SNS → SQS Subscriptions ───────────────────────────────────────────────────
# RawMessageDelivery=true: consumers receive the event envelope directly,
# without the SNS notification wrapper.

echo "Wiring SNS → SQS subscriptions..."

# Helper: subscribe a queue to a topic
subscribe() {
  local TOPIC_NAME="$1"
  local QUEUE_NAME="$2"

  TOPIC_ARN=$(awslocal sns list-topics \
    --query "Topics[?ends_with(TopicArn, ':${TOPIC_NAME}')].TopicArn | [0]" \
    --output text 2>/dev/null || true)

  QUEUE_URL=$(awslocal sqs get-queue-url \
    --queue-name "$QUEUE_NAME" \
    --query QueueUrl --output text 2>/dev/null || true)

  QUEUE_ARN=$(awslocal sqs get-queue-attributes \
    --queue-url "$QUEUE_URL" \
    --attribute-names QueueArn \
    --query Attributes.QueueArn --output text 2>/dev/null || true)

  if [ -n "$TOPIC_ARN" ] && [ -n "$QUEUE_ARN" ]; then
    awslocal sns subscribe \
      --topic-arn "$TOPIC_ARN" \
      --protocol sqs \
      --notification-endpoint "$QUEUE_ARN" \
      --attributes RawMessageDelivery=true \
      2>/dev/null || echo "Subscription ${TOPIC_NAME} → ${QUEUE_NAME} already exists"
  else
    echo "WARNING: could not subscribe ${QUEUE_NAME} to ${TOPIC_NAME} (topic or queue not found)"
  fi
}

subscribe "komodo-order-events-local.fifo"   "komodo-order-events-inventory-local.fifo"
subscribe "komodo-order-events-local.fifo"   "komodo-order-events-communications-local.fifo"
subscribe "komodo-order-events-local.fifo"   "komodo-order-events-loyalty-local.fifo"
subscribe "komodo-payment-events-local.fifo" "komodo-payment-events-order-local.fifo"
subscribe "komodo-payment-events-local.fifo" "komodo-payment-events-communications-local.fifo"
subscribe "komodo-user-events-local.fifo"    "komodo-user-events-loyalty-local.fifo"
subscribe "komodo-user-events-local.fifo"    "komodo-user-events-communications-local.fifo"

echo "SQS/SNS initialized successfully"
