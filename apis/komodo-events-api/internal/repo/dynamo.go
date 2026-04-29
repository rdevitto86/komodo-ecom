package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"komodo-events-api/internal/relay"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoRepository persists event envelopes to DynamoDB before dispatch.
type DynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoRepository(client *dynamodb.Client, tableName string) *DynamoRepository {
	return &DynamoRepository{client: client, tableName: tableName}
}

// SaveEvent writes the envelope to DynamoDB with status=pending and a 7-day TTL.
// PK=EVENT#<id>, SK=DOMAIN#<domain>.
func (r *DynamoRepository) SaveEvent(ctx context.Context, env relay.EventEnvelope) error {
	payload, err := json.Marshal(env.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	domain := domainFromType(string(env.Type))
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]types.AttributeValue{
			"event_id":    &types.AttributeValueMemberS{Value: "EVENT#" + env.ID},
			"domain":      &types.AttributeValueMemberS{Value: "DOMAIN#" + domain},
			"event_type":  &types.AttributeValueMemberS{Value: string(env.Type)},
			"source":      &types.AttributeValueMemberS{Value: string(env.Source)},
			"version":     &types.AttributeValueMemberS{Value: env.Version},
			"occurred_at": &types.AttributeValueMemberS{Value: env.OccurredAt.UTC().Format(time.RFC3339)},
			"payload":     &types.AttributeValueMemberS{Value: string(payload)},
			"status":      &types.AttributeValueMemberS{Value: "pending"},
			"expires_at":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expiresAt)},
		},
	})
	if err != nil {
		return fmt.Errorf("put event %s: %w", env.ID, err)
	}
	return nil
}

// UpdateEventStatus sets the status field on an existing event record.
func (r *DynamoRepository) UpdateEventStatus(ctx context.Context, eventID, domain, status string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: "EVENT#" + eventID},
			"domain":   &types.AttributeValueMemberS{Value: "DOMAIN#" + domain},
		},
		UpdateExpression: aws.String("SET #s = :s"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		return fmt.Errorf("update event %s status: %w", eventID, err)
	}
	return nil
}

func domainFromType(eventType string) string {
	if idx := strings.IndexByte(eventType, '.'); idx > 0 {
		return eventType[:idx]
	}
	return eventType
}
