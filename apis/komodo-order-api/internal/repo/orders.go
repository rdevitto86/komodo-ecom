package repo

import (
	"context"
	"errors"
	"fmt"
	"os"

	"komodo-order-api/internal/config"

	"komodo-order-api/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamo"
)

// table is resolved at package init from the env (populated by secrets-manager bootstrap).
var table = os.Getenv(config.DYNAMODB_ORDERS_TABLE)

// orderRecord is the METADATA row for an order, keyed by PK=ORDER#<id>, SK=METADATA.
// GSI1 enables listing all orders for a user in creation-time order.
type orderRecord struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI1PK    string `dynamodbav:"GSI1PK"`
	GSI1SK    string `dynamodbav:"GSI1SK"`
	ID        string `dynamodbav:"id"`
	DisplayID string `dynamodbav:"display_id"`
	UserID    string `dynamodbav:"user_id"`
	Email     string `dynamodbav:"email"`
	Status    string `dynamodbav:"status"`
	CreatedAt string `dynamodbav:"created_at"`
	UpdatedAt string `dynamodbav:"updated_at"`

	// Nested sub-documents stored as DynamoDB maps.
	Items   []models.OrderItem   `dynamodbav:"items"`
	Address models.OrderAddress  `dynamodbav:"address"`
	Payment models.OrderPayment  `dynamodbav:"payment"`
	Totals  models.OrderTotals   `dynamodbav:"totals"`
}

func orderPK(orderID string) string    { return "ORDER#" + orderID }
func userGSI1PK(userID string) string  { return "USER#" + userID }
func orderGSI1SK(createdAt, orderID string) string {
	return "ORDER#" + createdAt + "#" + orderID
}

// CreateOrder writes a new order to DynamoDB with a condition preventing overwrites.
// Key scheme:
//
//	PK=ORDER#<orderID>  SK=METADATA
//	GSI1PK=USER#<userID>  GSI1SK=ORDER#<createdAt>#<orderID>
func CreateOrder(ctx context.Context, order *models.Order) error {
	rec := orderRecord{
		PK:        orderPK(order.ID),
		SK:        "METADATA",
		GSI1PK:    userGSI1PK(order.UserID),
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
	if err := dynamo.WriteItem(ctx, table, av, false, nil, condition); err != nil {
		return fmt.Errorf("repo.CreateOrder: put: %w", err)
	}
	return nil
}

// GetOrder fetches the METADATA row for a single order by ID.
func GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	return nil, errors.New("not implemented")
}

// ListOrdersByUser returns all orders for a user, sorted by creation time descending,
// via the GSI1 index (GSI1PK=USER#<userID>).
func ListOrdersByUser(ctx context.Context, userID string) ([]*models.Order, error) {
	return nil, errors.New("not implemented")
}

// recordToOrder converts a raw DynamoDB attribute map to the Order model.
// Used by GetOrder and ListOrdersByUser once implemented.
func recordToOrder(raw map[string]ddbTypes.AttributeValue) (*models.Order, error) {
	var rec orderRecord
	if err := attributevalue.UnmarshalMap(raw, &rec); err != nil {
		return nil, fmt.Errorf("repo: unmarshal order record: %w", err)
	}
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
	}, nil
}
