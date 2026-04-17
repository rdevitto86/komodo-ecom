package repo

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"komodo-cart-api/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamo"
)

// table is resolved once at package init from config (set by secrets-manager bootstrap).
var table = os.Getenv("DYNAMODB_CARTS_TABLE")

// cartMetaRecord is the METADATA row for a user's cart.
type cartMetaRecord struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	UserID    string `dynamodbav:"user_id"`
	UpdatedAt string `dynamodbav:"updated_at"`
}

// cartItemRecord is an ITEM#<itemId> row in the cart table.
type cartItemRecord struct {
	PK             string `dynamodbav:"PK"`
	SK             string `dynamodbav:"SK"`
	ItemID         string `dynamodbav:"item_id"`
	SKU            string `dynamodbav:"sku"`
	Name           string `dynamodbav:"name"`
	Quantity       int    `dynamodbav:"quantity"`
	UnitPriceCents int    `dynamodbav:"unit_price_cents"`
	ImageURL       string `dynamodbav:"image_url"`
	UpdatedAt      string `dynamodbav:"updated_at"`
}

func cartPK(userID string) string { return "CART#" + userID }
func itemSK(itemID string) string { return "ITEM#" + itemID }

// cartIDFor returns a deterministic UUID for a user's cart.
// Computed on the fly — never stored in DynamoDB.
func cartIDFor(userID string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(userID)).String()
}

// GetCart fetches all rows for the user's cart and assembles the Cart model.
// Returns an empty Cart (not nil) with zero items if no rows exist.
func GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	rows, err := dynamo.QueryAll(ctx, dynamo.QueryInput{
		TableName:              table,
		KeyConditionExpression: "PK = :pk",
		ExpressionValues: map[string]ddbTypes.AttributeValue{
			":pk": &ddbTypes.AttributeValueMemberS{Value: cartPK(userID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetCart: query: %w", err)
	}

	cart := &models.Cart{
		ID:        cartIDFor(userID),
		UserID:    userID,
		Items:     []models.CartItem{},
		UpdatedAt: time.Now().UTC(),
	}

	for _, raw := range rows {
		skAttr, ok := raw["SK"]
		if !ok {
			continue
		}
		skMember, ok := skAttr.(*ddbTypes.AttributeValueMemberS)
		if !ok {
			continue
		}
		sk := skMember.Value

		switch {
		case sk == "METADATA":
			var meta cartMetaRecord
			if err := attributevalue.UnmarshalMap(raw, &meta); err != nil {
				// Non-fatal: metadata row malformed; keep default UpdatedAt.
				continue
			}
			if t, err := time.Parse(time.RFC3339, meta.UpdatedAt); err == nil {
				cart.UpdatedAt = t
			}

		case strings.HasPrefix(sk, "ITEM#"):
			var rec cartItemRecord
			if err := attributevalue.UnmarshalMap(raw, &rec); err != nil {
				// Skip malformed item rows rather than failing the whole cart read.
				continue
			}
			cart.Items = append(cart.Items, models.CartItem{
				ItemID:         rec.ItemID,
				SKU:            rec.SKU,
				Name:           rec.Name,
				Quantity:       rec.Quantity,
				UnitPriceCents: rec.UnitPriceCents,
				ImageURL:       rec.ImageURL,
			})
		}
	}

	return cart, nil
}

// PutCartItem upserts an ITEM#<itemID> row and refreshes the METADATA row.
func PutCartItem(ctx context.Context, userID string, item models.CartItem) error {
	now := time.Now().UTC().Format(time.RFC3339)
	pk := cartPK(userID)

	rec := cartItemRecord{
		PK:             pk,
		SK:             itemSK(item.ItemID),
		ItemID:         item.ItemID,
		SKU:            item.SKU,
		Name:           item.Name,
		Quantity:       item.Quantity,
		UnitPriceCents: item.UnitPriceCents,
		ImageURL:       item.ImageURL,
		UpdatedAt:      now,
	}
	if err := dynamo.WriteItemFrom(ctx, table, rec, false, nil, nil); err != nil {
		return fmt.Errorf("repo.PutCartItem: write item: %w", err)
	}

	meta := cartMetaRecord{
		PK:        pk,
		SK:        "METADATA",
		UserID:    userID,
		UpdatedAt: now,
	}
	if err := dynamo.WriteItemFrom(ctx, table, meta, false, nil, nil); err != nil {
		return fmt.Errorf("repo.PutCartItem: write metadata: %w", err)
	}
	return nil
}

// UpdateCartItemQuantity updates the quantity and updated_at of an existing item row.
func UpdateCartItemQuantity(ctx context.Context, userID, itemID string, qty int) error {
	now := time.Now().UTC().Format(time.RFC3339)

	key, err := dynamo.BuildKey("PK", cartPK(userID), "SK", itemSK(itemID))
	if err != nil {
		return fmt.Errorf("repo.UpdateCartItemQuantity: build key: %w", err)
	}

	exprValues := map[string]ddbTypes.AttributeValue{
		":qty": &ddbTypes.AttributeValueMemberN{Value: fmt.Sprintf("%d", qty)},
		":ts":  &ddbTypes.AttributeValueMemberS{Value: now},
	}

	_, err = dynamo.UpdateItem(ctx, table, key,
		"SET quantity = :qty, updated_at = :ts",
		exprValues, nil, nil)
	if err != nil {
		return fmt.Errorf("repo.UpdateCartItemQuantity: update: %w", err)
	}
	return nil
}

// DeleteCartItem removes a single item row from the cart table.
func DeleteCartItem(ctx context.Context, userID, itemID string) error {
	key, err := dynamo.BuildKey("PK", cartPK(userID), "SK", itemSK(itemID))
	if err != nil {
		return fmt.Errorf("repo.DeleteCartItem: build key: %w", err)
	}
	if err := dynamo.DeleteItem(ctx, table, key, false, nil, nil); err != nil {
		return fmt.Errorf("repo.DeleteCartItem: delete: %w", err)
	}
	return nil
}

// ClearCart removes all ITEM#* rows and resets the METADATA updated_at.
func ClearCart(ctx context.Context, userID string) error {
	pk := cartPK(userID)

	rows, err := dynamo.QueryAll(ctx, dynamo.QueryInput{
		TableName:              table,
		KeyConditionExpression: "PK = :pk AND begins_with(SK, :prefix)",
		ExpressionValues: map[string]ddbTypes.AttributeValue{
			":pk":     &ddbTypes.AttributeValueMemberS{Value: pk},
			":prefix": &ddbTypes.AttributeValueMemberS{Value: "ITEM#"},
		},
	})
	if err != nil {
		return fmt.Errorf("repo.ClearCart: query items: %w", err)
	}

	if len(rows) > 0 {
		keys := make([]map[string]ddbTypes.AttributeValue, 0, len(rows))
		for _, row := range rows {
			pkAttr, hasPK := row["PK"]
			skAttr, hasSK := row["SK"]
			if !hasPK || !hasSK {
				continue
			}
			keys = append(keys, map[string]ddbTypes.AttributeValue{"PK": pkAttr, "SK": skAttr})
		}
		if len(keys) > 0 {
			if err := dynamo.DeleteItem(ctx, table, nil, true, keys, nil); err != nil {
				return fmt.Errorf("repo.ClearCart: batch delete: %w", err)
			}
		}
	}

	// Reset METADATA updated_at to signal the cart was cleared.
	meta := cartMetaRecord{
		PK:        pk,
		SK:        "METADATA",
		UserID:    userID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := dynamo.WriteItemFrom(ctx, table, meta, false, nil, nil); err != nil {
		return fmt.Errorf("repo.ClearCart: write metadata: %w", err)
	}
	return nil
}

// ItemExists returns true if an ITEM#<itemID> row exists for the user's cart.
func ItemExists(ctx context.Context, userID, itemID string) (bool, error) {
	key, err := dynamo.BuildKey("PK", cartPK(userID), "SK", itemSK(itemID))
	if err != nil {
		return false, fmt.Errorf("repo.ItemExists: build key: %w", err)
	}

	var rec cartItemRecord
	err = dynamo.GetItemAs(ctx, table, key, false, nil, &rec)
	if err != nil {
		// SDK returns an error for item-not-found — treat that as not-found (not an error).
		return false, nil
	}
	return rec.ItemID != "", nil
}
