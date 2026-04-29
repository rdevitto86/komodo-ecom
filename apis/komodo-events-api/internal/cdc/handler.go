package cdc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"komodo-events-api/internal/config"

	komodoEvents "github.com/rdevitto86/komodo-forge-sdk-go/events"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-lambda-go/events"
)

// Handler processes DynamoDB Stream records and forwards business events
// to the event-bus-api /internal/dispatch endpoint via HTTP POST.
// One handler instance is shared across all Lambda invocations.
type Handler struct {
	eventBusURL string
	httpClient  *http.Client
	env         string
}

func NewHandler(eventBusURL string) *Handler {
	return &Handler{
		eventBusURL: eventBusURL,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		env:         os.Getenv(config.ENV),
	}
}

// Handle is the Lambda entry point for DynamoDB Stream events.
// Per-record failures are logged but do not fail the overall invocation.
//
// TODO: emit a CloudWatch metric (or a structured log entry with a fixed key)
// on per-record dispatch failure so an alarm can fire before the DLQ fills.
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
		return nil
	}

	evt := komodoEvents.New(result.EventType, result.Source, result.EntityID, result.EntityType, result.Payload)
	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event envelope: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		h.eventBusURL+"/internal/dispatch", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build dispatch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("dispatch http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("event-bus returned %d for event %s", resp.StatusCode, evt.ID)
	}

	logger.Info("event dispatched",
		logger.Attr("event_id", evt.ID),
		logger.Attr("event_type", string(evt.Type)),
		logger.Attr("entity_id", evt.EntityID),
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
