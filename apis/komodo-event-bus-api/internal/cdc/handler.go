package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"komodo-forge-sdk-go/config"
	komodoEvents "komodo-forge-sdk-go/events"
	logger "komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// Handler processes DynamoDB Stream records and publishes business events
// to SNS FIFO topics. One handler instance is shared across all invocations.
type Handler struct {
	sns            *sns.Client
	topicARNPrefix string // e.g. "arn:aws:sns:us-east-1:123456789012:komodo-"
	env            string // e.g. "prod", "staging", "dev"
}

func NewHandler(snsClient *sns.Client, topicARNPrefix string) *Handler {
	return &Handler{
		sns:            snsClient,
		topicARNPrefix: topicARNPrefix,
		env:            config.GetConfigValue("ENV"),
	}
}

// Handle is the Lambda entry point for DynamoDB Stream events.
// Errors from individual records are logged but do not fail the batch — a
// single bad record should not stall the stream. Genuinely unprocessable
// messages are caught by the SQS DLQ on the consumer side.
//
// TODO: emit a CloudWatch metric (or a structured log entry with a fixed key)
// on per-record publish failure so an alarm can fire before the DLQ fills.
//
// TODO: consider sns.PublishBatch when len(event.Records) > 1 — reduces SNS
// API calls and latency. Batch limit: 10 messages / 256 KB total.
func (h *Handler) Handle(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		if err := h.processRecord(ctx, record); err != nil {
			logger.Error("failed to process stream record", err,
				logger.Attr("event_id", record.EventID),
				logger.Attr("event_name", record.EventName),
			)
		}
	}
	return nil
}

func (h *Handler) processRecord(ctx context.Context, record events.DynamoDBEventRecord) error {
	tableName := tableNameFromARN(record.EventSourceArn)
	if tableName == "" {
		return fmt.Errorf("could not extract table name from stream ARN: %s", record.EventSourceArn)
	}

	result, ok := classify(tableName, record.EventName, record.Change.OldImage, record.Change.NewImage)
	if !ok {
		return nil // not a business-meaningful change — skip silently
	}

	evt := komodoEvents.New(result.EventType, result.Source, result.EntityID, result.EntityType, result.Payload)

	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event envelope: %w", err)
	}

	// SNS FIFO: MessageGroupId preserves per-entity ordering while allowing
	// cross-entity parallelism. MessageDeduplicationId prevents duplicates.
	topicARN := fmt.Sprintf("%s%s-events-%s.fifo", h.topicARNPrefix, result.Domain, h.env)
	_, err = h.sns.Publish(ctx, &sns.PublishInput{
		TopicArn:               aws.String(topicARN),
		Message:                aws.String(string(body)),
		MessageGroupId:         aws.String(evt.EntityID),
		MessageDeduplicationId: aws.String(evt.ID),
		MessageAttributes: map[string]snsTypes.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(evt.Type)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("sns publish to %s: %w", topicARN, err)
	}

	logger.Info("event published",
		logger.Attr("event_id", evt.ID),
		logger.Attr("event_type", string(evt.Type)),
		logger.Attr("entity_id", evt.EntityID),
		logger.Attr("topic_arn", topicARN),
	)
	return nil
}

// tableNameFromARN extracts the DynamoDB table name from a stream event source ARN.
// Format: arn:aws:dynamodb:REGION:ACCOUNT:table/TABLE_NAME/stream/TIMESTAMP
func tableNameFromARN(arn string) string {
	parts := strings.SplitN(arn, "/", 3)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
