package repo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"komodo-order-api/internal/config"
	"komodo-order-api/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dyn "github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamo"
)

// table is resolved at package init from the env (populated by secrets-manager bootstrap).
var table = os.Getenv(config.DYNAMODB_ORDERS_TABLE)

// gsi1Name is the GSI used for user-scoped order listing.
const gsi1Name = "GSI1"

// rawClient holds the raw DynamoDB client used for GSI queries (which require
// DynamoDB-native API calls not covered by the forge SDK's dynamo.QueryAll wrapper).
// Initialised lazily on first use via rawClientOnce.
var (
	rawClient     *dynamodb.Client
	rawClientOnce sync.Once
)

// getClient returns the raw DynamoDB client, initializing it on the first call.
// Uses the same env vars as dynamo.Init in main.go.
func getClient() (*dynamodb.Client, error) {
	var initErr error
	rawClientOnce.Do(func() {
		region := os.Getenv(config.AWS_REGION)
		endpoint := os.Getenv(config.DYNAMODB_ENDPOINT)
		accessKey := os.Getenv(config.DYNAMODB_ACCESS_KEY)
		secretKey := os.Getenv(config.DYNAMODB_SECRET_KEY)

		opts := []func(*awsCfg.LoadOptions) error{
			awsCfg.WithRegion(region),
		}
		if accessKey != "" && secretKey != "" {
			opts = append(opts, awsCfg.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
			))
		}

		cfg, err := awsCfg.LoadDefaultConfig(context.Background(), opts...)
		if err != nil {
			initErr = fmt.Errorf("repo.getClient: load aws config: %w", err)
			return
		}

		clientOpts := []func(*dynamodb.Options){}
		if endpoint != "" {
			clientOpts = append(clientOpts, func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(endpoint)
			})
		}

		rawClient = dynamodb.NewFromConfig(cfg, clientOpts...)
	})
	if initErr != nil {
		return nil, initErr
	}
	return rawClient, nil
}

// orderRecord is the METADATA row for an order, keyed by PK=ORDER#<id>, SK=METADATA.
// GSI1 enables listing all orders for a user in creation-time order.
// GSI1PK is USER#<userId> for registered accounts and GUEST#<uuid> for guests.
type orderRecord struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI1PK    string `dynamodbav:"GSI1PK"`
	GSI1SK    string `dynamodbav:"GSI1SK"`
	ID        string `dynamodbav:"id"`
	DisplayID string `dynamodbav:"display_id"`
	UserID    string `dynamodbav:"user_id"`
	Email     string `dynamodbav:"email,omitempty"`
	Status    string `dynamodbav:"status"`
	CreatedAt string `dynamodbav:"created_at"`
	UpdatedAt string `dynamodbav:"updated_at"`

	// Nested sub-documents stored as DynamoDB maps.
	Items   []models.OrderItem  `dynamodbav:"items"`
	Address models.OrderAddress `dynamodbav:"address"`
	Payment models.OrderPayment `dynamodbav:"payment"`
	Totals  models.OrderTotals  `dynamodbav:"totals"`
}

func orderPK(orderID string) string   { return "ORDER#" + orderID }
func userGSI1PK(userID string) string { return "USER#" + userID }
func orderGSI1SK(createdAt, orderID string) string {
	return "ORDER#" + createdAt + "#" + orderID
}

// CreateOrder writes a new order to DynamoDB with a condition preventing overwrites.
// Key scheme:
//
//	PK=ORDER#<orderID>  SK=METADATA
//	GSI1PK=<order.UserID>  GSI1SK=ORDER#<createdAt>#<orderID>
//
// order.UserID is expected to already carry the key prefix: USER#<id> for
// registered accounts, GUEST#<uuid> for unauthenticated placements.
func CreateOrder(ctx context.Context, order *models.Order) error {
	rec := orderRecord{
		PK:        orderPK(order.ID),
		SK:        "METADATA",
		GSI1PK:    order.UserID, // already prefixed: USER#<id> or GUEST#<uuid>
		GSI1SK:    orderGSI1SK(order.CreatedAt, order.ID),
		ID:        order.ID,
		DisplayID: order.DisplayID,
		UserID:    order.UserID,
		Email:     order.Email,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Items:     order.Items,
		Address:   order.Address,
		Payment:   order.Payment,
		Totals:    order.Totals,
	}

	av, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("repo.CreateOrder: marshal: %w", err)
	}

	// attribute_not_exists(PK) prevents double-writes for the same orderID.
	condition := aws.String("attribute_not_exists(PK)")
	if err := dyn.WriteItem(ctx, table, av, false, nil, condition); err != nil {
		return fmt.Errorf("repo.CreateOrder: put: %w", err)
	}
	return nil
}

// GetOrder fetches the METADATA row for a single order by its primary key.
// Returns models.ErrNotFound (wrapped) when the item does not exist in DynamoDB.
func GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	key, err := dyn.BuildKey("PK", orderPK(orderID), "SK", "METADATA")
	if err != nil {
		return nil, fmt.Errorf("repo.GetOrder: build key: %w", err)
	}

	var rec orderRecord
	if err := dyn.GetItemAs(ctx, table, key, false, nil, &rec); err != nil {
		// GetItemAs returns an error for item-not-found — surface as ErrNotFound.
		return nil, fmt.Errorf("repo.GetOrder: %w", models.ErrNotFound)
	}
	if rec.ID == "" {
		// Item key existed but record has no data — treat as not found.
		return nil, fmt.Errorf("repo.GetOrder: empty record: %w", models.ErrNotFound)
	}

	return recordToModel(&rec), nil
}

// ListOrdersByUser returns a page of orders for a user via the GSI1 index
// (GSI1PK=USER#<userID>), sorted newest-first (descending by GSI1SK).
//
// limit controls the page size (0 defaults to 20; max is capped at 100).
// cursor is an opaque DynamoDB continuation token from a previous call (empty = first page).
// Returns the orders, a next-page cursor (empty string if no more pages), and any error.
func ListOrdersByUser(ctx context.Context, userID string, limit int, cursor string) ([]*models.Order, string, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	client, err := getClient()
	if err != nil {
		return nil, "", fmt.Errorf("repo.ListOrdersByUser: init client: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(table),
		IndexName:              aws.String(gsi1Name),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":pk": &ddbTypes.AttributeValueMemberS{Value: userGSI1PK(userID)},
		},
		ScanIndexForward: aws.Bool(false), // newest-first
		Limit:            aws.Int32(int32(limit)),
	}

	if cursor != "" {
		startKey, decodeErr := decodeCursor(cursor)
		if decodeErr != nil {
			return nil, "", fmt.Errorf("repo.ListOrdersByUser: decode cursor: %w", decodeErr)
		}
		input.ExclusiveStartKey = startKey
	}

	out, err := client.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("repo.ListOrdersByUser: query: %w", err)
	}

	orders := make([]*models.Order, 0, len(out.Items))
	for _, raw := range out.Items {
		var rec orderRecord
		if err := attributevalue.UnmarshalMap(raw, &rec); err != nil {
			// Skip malformed rows — non-fatal.
			continue
		}
		orders = append(orders, recordToModel(&rec))
	}

	nextCursor := ""
	if len(out.LastEvaluatedKey) > 0 {
		nextCursor, err = encodeCursor(out.LastEvaluatedKey)
		if err != nil {
			return nil, "", fmt.Errorf("repo.ListOrdersByUser: encode cursor: %w", err)
		}
	}

	return orders, nextCursor, nil
}

// UpdateOrderStatus conditionally sets the order status to newStatus.
// The write is rejected if the stored status != expectedStatus, guarding against
// concurrent writes racing to transition the same order.
//
// "status" is a DynamoDB reserved word — #s is used as the expression attribute name.
func UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus, expectedStatus models.OrderStatus) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("repo.UpdateOrderStatus: init client: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]ddbTypes.AttributeValue{
			"PK": &ddbTypes.AttributeValueMemberS{Value: orderPK(orderID)},
			"SK": &ddbTypes.AttributeValueMemberS{Value: "METADATA"},
		},
		UpdateExpression:    aws.String("SET #s = :new, updated_at = :ua"),
		ConditionExpression: aws.String("attribute_exists(PK) AND #s = :expected"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":new":      &ddbTypes.AttributeValueMemberS{Value: string(newStatus)},
			":expected": &ddbTypes.AttributeValueMemberS{Value: string(expectedStatus)},
			":ua":       &ddbTypes.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
	}

	_, err = client.UpdateItem(ctx, input)
	if err != nil {
		var condErr *ddbTypes.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return fmt.Errorf("repo.UpdateOrderStatus: condition failed: %w", models.ErrInvalidTransition)
		}
		return fmt.Errorf("repo.UpdateOrderStatus: update: %w", err)
	}
	return nil
}

// recordToModel converts an unmarshaled orderRecord to the Order domain model.
func recordToModel(rec *orderRecord) *models.Order {
	return &models.Order{
		ID:        rec.ID,
		DisplayID: rec.DisplayID,
		UserID:    rec.UserID,
		Email:     rec.Email,
		Status:    models.OrderStatus(rec.Status),
		Items:     rec.Items,
		Address:   rec.Address,
		Payment:   rec.Payment,
		Totals:    rec.Totals,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

// encodeCursor serialises a DynamoDB LastEvaluatedKey to a URL-safe base64 JSON string.
// The GSI1 pagination key consists of PK, SK, GSI1PK, GSI1SK — all string attributes.
func encodeCursor(key map[string]ddbTypes.AttributeValue) (string, error) {
	simple := make(map[string]string, len(key))
	for k, v := range key {
		sv, ok := v.(*ddbTypes.AttributeValueMemberS)
		if !ok {
			return "", fmt.Errorf("encodeCursor: unexpected non-string attribute %q in pagination key", k)
		}
		simple[k] = sv.Value
	}
	b, err := json.Marshal(simple)
	if err != nil {
		return "", fmt.Errorf("encodeCursor: marshal: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// decodeCursor reverses encodeCursor.
func decodeCursor(cursor string) (map[string]ddbTypes.AttributeValue, error) {
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("decodeCursor: base64: %w", err)
	}
	var simple map[string]string
	if err := json.Unmarshal(b, &simple); err != nil {
		return nil, fmt.Errorf("decodeCursor: unmarshal: %w", err)
	}
	result := make(map[string]ddbTypes.AttributeValue, len(simple))
	for k, v := range simple {
		result[k] = &ddbTypes.AttributeValueMemberS{Value: v}
	}
	return result, nil
}

