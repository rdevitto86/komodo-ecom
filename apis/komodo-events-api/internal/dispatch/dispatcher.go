package dispatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"komodo-events-api/internal/relay"
	"komodo-events-api/internal/repo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// Dispatcher fans out a validated event envelope to all active HTTP subscribers.
type Dispatcher struct {
	dynamo             *dynamodb.Client
	eventsTable        string
	subscriptionsTable string
	httpClient         *http.Client
	repo               *repo.DynamoRepository
}

func NewDispatcher(dynamo *dynamodb.Client, eventsTable, subscriptionsTable string) *Dispatcher {
	return &Dispatcher{
		dynamo:             dynamo,
		eventsTable:        eventsTable,
		subscriptionsTable: subscriptionsTable,
		httpClient:         &http.Client{Timeout: 5 * time.Second},
		repo:               repo.NewDynamoRepository(dynamo, eventsTable),
	}
}

// Dispatch queries active subscribers for the event type and HTTP-POSTs the envelope to each.
// Partial subscriber failures are logged but do not fail the overall dispatch.
// The event status in komodo-events is updated to "dispatched" or "failed" after all calls.
func (d *Dispatcher) Dispatch(ctx context.Context, env relay.EventEnvelope) error {
	subscribers, err := d.querySubscribers(ctx, string(env.Type))
	if err != nil {
		return fmt.Errorf("query subscribers for %s: %w", env.Type, err)
	}

	if len(subscribers) == 0 {
		logger.Info("no active subscribers",
			logger.Attr("event_id", env.ID),
			logger.Attr("event_type", string(env.Type)),
		)
		d.updateStatus(ctx, env, "dispatched")
		return nil
	}

	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}

	anySuccess := false
	for _, url := range subscribers {
		if postErr := d.post(ctx, url, body, env.ID, string(env.Type)); postErr != nil {
			logger.Error("subscriber dispatch failed", postErr,
				logger.Attr("event_id", env.ID),
				logger.Attr("event_type", string(env.Type)),
				logger.Attr("subscriber_url", url),
			)
		} else {
			anySuccess = true
			logger.Info("event dispatched",
				logger.Attr("event_id", env.ID),
				logger.Attr("event_type", string(env.Type)),
				logger.Attr("subscriber_url", url),
			)
		}
	}

	status := "dispatched"
	if !anySuccess {
		status = "failed"
	}
	d.updateStatus(ctx, env, status)
	return nil
}

// HandleDispatch is the HTTP handler for POST /internal/dispatch.
// Called by the CDC Lambda; no auth required (VPC-internal).
func (d *Dispatcher) HandleDispatch(w http.ResponseWriter, r *http.Request) {
	var env relay.EventEnvelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := d.Dispatch(r.Context(), env); err != nil {
		logger.Error("internal dispatch failed", err,
			logger.Attr("event_id", env.ID),
			logger.Attr("event_type", string(env.Type)),
		)
		http.Error(w, "dispatch failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// querySubscribers returns the subscriber_url values for active subscriptions to the given event type.
func (d *Dispatcher) querySubscribers(ctx context.Context, eventType string) ([]string, error) {
	out, err := d.dynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.subscriptionsTable),
		KeyConditionExpression: aws.String("event_type = :et"),
		FilterExpression:       aws.String("active = :a"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":et": &types.AttributeValueMemberS{Value: eventType},
			":a":  &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	urls := make([]string, 0, len(out.Items))
	for _, item := range out.Items {
		if v, ok := item["subscriber_url"].(*types.AttributeValueMemberS); ok && v.Value != "" {
			urls = append(urls, v.Value)
		}
	}
	return urls, nil
}

func (d *Dispatcher) post(ctx context.Context, url string, body []byte, eventID, eventType string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-ID", eventID)
	req.Header.Set("X-Event-Type", eventType)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("subscriber returned %d", resp.StatusCode)
	}
	return nil
}

func (d *Dispatcher) updateStatus(ctx context.Context, env relay.EventEnvelope, status string) {
	domain := domainFromType(string(env.Type))
	if err := d.repo.UpdateEventStatus(ctx, env.ID, domain, status); err != nil {
		logger.Error("failed to update event status", err,
			logger.Attr("event_id", env.ID),
			logger.Attr("status", status),
		)
	}
}

func domainFromType(eventType string) string {
	if idx := strings.IndexByte(eventType, '.'); idx > 0 {
		return eventType[:idx]
	}
	return eventType
}
