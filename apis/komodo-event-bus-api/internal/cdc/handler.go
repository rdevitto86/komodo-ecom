package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	komodoEvents "github.com/rdevitto86/komodo-forge-sdk-go/events"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

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
		env:            os.Getenv("ENV"),
	}
}

// pendingPublish holds everything needed to publish one event to SNS.
type pendingPublish struct {
	topicARN  string
	message   string
	groupID   string // MessageGroupId — entityID for per-entity ordering
	dedupID   string // MessageDeduplicationId — event UUID
	eventID   string
	eventType string
}

// Handle is the Lambda entry point for DynamoDB Stream events. When the batch
// contains more than one record, messages are grouped by topic and sent via
// sns.PublishBatch (up to 10 per API call) to reduce SNS round-trips.
// Per-record failures are logged but do not fail the overall invocation.
//
// TODO: emit a CloudWatch metric (or a structured log entry with a fixed key)
// on per-record publish failure so an alarm can fire before the DLQ fills.
func (h *Handler) Handle(ctx context.Context, event events.DynamoDBEvent) error {
	if len(event.Records) > 1 {
		return h.handleBatch(ctx, event.Records)
	}
	if len(event.Records) == 1 {
		if err := h.processRecord(ctx, event.Records[0]); err != nil {
			logger.Error("failed to process stream record", err,
				logger.Attr("event_id", event.Records[0].EventID),
				logger.Attr("event_name", event.Records[0].EventName),
			)
		}
	}
	return nil
}

// toPendingPublish classifies a stream record and marshals the event envelope.
// Returns (zero, false, nil) for records that are not business-meaningful.
func (h *Handler) toPendingPublish(record events.DynamoDBEventRecord) (pendingPublish, bool, error) {
	tableName := tableNameFromARN(record.EventSourceArn)
	if tableName == "" {
		return pendingPublish{}, false, fmt.Errorf("could not extract table name from stream ARN: %s", record.EventSourceArn)
	}

	result, ok := classify(tableName, record.EventName, record.Change.OldImage, record.Change.NewImage)
	if !ok {
		return pendingPublish{}, false, nil
	}

	evt := komodoEvents.New(result.EventType, result.Source, result.EntityID, result.EntityType, result.Payload)
	body, err := json.Marshal(evt)
	if err != nil {
		return pendingPublish{}, false, fmt.Errorf("marshal event envelope: %w", err)
	}

	return pendingPublish{
		topicARN:  fmt.Sprintf("%s%s-events-%s.fifo", h.topicARNPrefix, result.Domain, h.env),
		message:   string(body),
		groupID:   evt.EntityID,
		dedupID:   evt.ID,
		eventID:   evt.ID,
		eventType: string(evt.Type),
	}, true, nil
}

func (h *Handler) processRecord(ctx context.Context, record events.DynamoDBEventRecord) error {
	pp, ok, err := h.toPendingPublish(record)
	if err != nil { return err }
	if !ok { return nil }

	// SNS FIFO: MessageGroupId preserves per-entity ordering while allowing
	// cross-entity parallelism. MessageDeduplicationId prevents duplicates.
	_, err = h.sns.Publish(ctx, &sns.PublishInput{
		TopicArn:               aws.String(pp.topicARN),
		Message:                aws.String(pp.message),
		MessageGroupId:         aws.String(pp.groupID),
		MessageDeduplicationId: aws.String(pp.dedupID),
		MessageAttributes: map[string]snsTypes.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(pp.eventType),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("sns publish to %s: %w", pp.topicARN, err)
	}

	logger.Info("event published",
		logger.Attr("event_id", pp.eventID),
		logger.Attr("event_type", pp.eventType),
		logger.Attr("entity_id", pp.groupID),
		logger.Attr("topic_arn", pp.topicARN),
	)
	return nil
}

// handleBatch classifies all records, groups them by SNS topic, then sends
// batches of up to 10 per topic via sns.PublishBatch.
func (h *Handler) handleBatch(ctx context.Context, records []events.DynamoDBEventRecord) error {
	byTopic := map[string][]pendingPublish{}
	for _, record := range records {
		pp, ok, err := h.toPendingPublish(record)
		if err != nil {
			logger.Error("failed to classify stream record", err,
				logger.Attr("event_id", record.EventID),
				logger.Attr("event_name", record.EventName),
			)
			continue
		}
		if !ok {
			continue
		}
		byTopic[pp.topicARN] = append(byTopic[pp.topicARN], pp)
	}

	for topicARN, pending := range byTopic {
		for i := 0; i < len(pending); i += 10 {
			end := i + 10
			if end > len(pending) {
				end = len(pending)
			}
			h.publishBatch(ctx, topicARN, pending[i:end])
		}
	}
	return nil
}

// publishBatch sends up to 10 pre-classified events to a single SNS FIFO topic.
// Partial failures are logged individually; the batch is not retried.
func (h *Handler) publishBatch(ctx context.Context, topicARN string, pending []pendingPublish) {
	entries := make([]snsTypes.PublishBatchRequestEntry, len(pending))
	for i, pp := range pending {
		entries[i] = snsTypes.PublishBatchRequestEntry{
			Id:                     aws.String(strconv.Itoa(i)),
			Message:                aws.String(pp.message),
			MessageGroupId:         aws.String(pp.groupID),
			MessageDeduplicationId: aws.String(pp.dedupID),
			MessageAttributes: map[string]snsTypes.MessageAttributeValue{
				"event_type": {
					DataType:    aws.String("String"),
					StringValue: aws.String(pp.eventType),
				},
			},
		}
	}

	out, err := h.sns.PublishBatch(ctx, &sns.PublishBatchInput{
		TopicArn:                   aws.String(topicARN),
		PublishBatchRequestEntries: entries,
	})
	if err != nil {
		// Entire batch failed — log each entry individually for visibility.
		for _, pp := range pending {
			logger.Error("batch publish failed", err,
				logger.Attr("event_id", pp.eventID),
				logger.Attr("event_type", pp.eventType),
				logger.Attr("topic_arn", topicARN),
			)
		}
		return
	}

	for _, f := range out.Failed {
		idx, _ := strconv.Atoi(aws.ToString(f.Id))
		pp := pending[idx]
		logger.Error("sns rejected event in batch",
			fmt.Errorf("code=%s: %s", aws.ToString(f.Code), aws.ToString(f.Message)),
			logger.Attr("event_id", pp.eventID),
			logger.Attr("event_type", pp.eventType),
			logger.Attr("sender_fault", f.SenderFault),
		)
	}

	for _, s := range out.Successful {
		idx, _ := strconv.Atoi(aws.ToString(s.Id))
		pp := pending[idx]
		logger.Info("event published",
			logger.Attr("event_id", pp.eventID),
			logger.Attr("event_type", pp.eventType),
			logger.Attr("entity_id", pp.groupID),
			logger.Attr("message_id", aws.ToString(s.MessageId)),
		)
	}
}

// tableNameFromARN extracts the DynamoDB table name from a stream event source ARN.
// Format: arn:aws:dynamodb:REGION:ACCOUNT:table/TABLE_NAME/stream/TIMESTAMP
func tableNameFromARN(arn string) string {
	parts := strings.SplitN(arn, "/", 3)
	if len(parts) < 2 { return "" }
	return parts[1]
}
