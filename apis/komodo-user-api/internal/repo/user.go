package repo

import (
	"context"
	"fmt"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamodb"
	"github.com/rdevitto86/komodo-forge-sdk-go/config"
	"komodo-user-api/internal/models"
)

// table is resolved at startup from config.
// Set DYNAMODB_TABLE=komodo-users for LocalStack and production alike —
// the forge DynamoDB client switches endpoints transparently via DYNAMODB_ENDPOINT.
var table = config.GetConfigValue("DYNAMODB_TABLE")

// TODO: define full key schema in data-model.md before finalizing these functions.
// Placeholder pattern used below: PK=USER#<id>, SK=PROFILE|ADDR#<id>|PREFS

// userRecord is the internal DynamoDB representation of a user.
// It includes sensitive fields (password_hash, PK, SK) that must never be
// serialized into an HTTP response. GetUser maps this to models.User before returning.
type userRecord struct {
	PK           string `dynamodbav:"PK"`
	SK           string `dynamodbav:"SK"`
	UserID       string `dynamodbav:"user_id"`
	Username     string `dynamodbav:"username"`
	Email        string `dynamodbav:"email"`
	Phone        string `dynamodbav:"phone"`
	FirstName    string `dynamodbav:"first_name"`
	MiddleInitial string `dynamodbav:"middle_initial"`
	LastName     string `dynamodbav:"last_name"`
	AvatarURL    string `dynamodbav:"avatar_url"`
	PasswordHash string `dynamodbav:"password_hash"` // never leaves this package
}

func (r *userRecord) toModel() *models.User {
	return &models.User{
		UserID:        r.UserID,
		Email:         r.Email,
		Phone:         r.Phone,
		FirstName:     r.FirstName,
		MiddleInitial: r.MiddleInitial,
		LastName:      r.LastName,
		AvatarURL:     r.AvatarURL,
	}
}

// GetUser retrieves a user's profile record.
// Reads into the internal userRecord (which includes password_hash) and maps
// to the public models.User — password_hash never escapes this function.
func GetUser(ctx context.Context, userID string) (*models.User, error) {
	key, err := dynamodb.BuildKey("PK", "USER#"+userID, "SK", "PROFILE")
	if err != nil {
		return nil, fmt.Errorf("repo.GetUser: build key: %w", err)
	}

	var record userRecord
	if err := dynamodb.GetItemAs(ctx, table, key, false, nil, &record); err != nil {
		return nil, fmt.Errorf("repo.GetUser: %w", err)
	}
	return record.toModel(), nil
}

// CreateUser writes a new user profile record.
// TODO: marshal PK/SK onto the item before writing once key schema is finalized.
func CreateUser(ctx context.Context, user *models.User) error {
	if err := dynamodb.WriteItemFrom(ctx, table, user, false, nil, nil); err != nil {
		return fmt.Errorf("repo.CreateUser: %w", err)
	}
	return nil
}

// UpdateUser applies a partial update to a user profile record.
// TODO: build update expression dynamically from changed fields once schema is settled.
func UpdateUser(ctx context.Context, userID string, _ *models.User) (*models.User, error) {
	key, err := dynamodb.BuildKey("PK", "USER#"+userID, "SK", "PROFILE")
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateUser: build key: %w", err)
	}

	var record userRecord
	// TODO: replace placeholder expression with field-specific SET clauses.
	if err := dynamodb.UpdateItemAs(ctx, table, key, "SET #placeholder = :placeholder", nil, nil, nil, &record); err != nil {
		return nil, fmt.Errorf("repo.UpdateUser: %w", err)
	}
	return record.toModel(), nil
}

// DeleteUser removes all items belonging to a user (profile, addresses, preferences).
// Queries all items with the user's partition key then batch deletes them in one pass,
// preventing orphaned records when a user account is closed.
// TODO: migrate to a DynamoDB transaction once write patterns are finalized.
func DeleteUser(ctx context.Context, userID string) error {
	items, err := dynamodb.QueryAll(ctx, dynamodb.QueryInput{
		TableName:              table,
		KeyConditionExpression: "PK = :pk",
		ExpressionValues: map[string]ddbTypes.AttributeValue{
			":pk": &ddbTypes.AttributeValueMemberS{Value: "USER#" + userID},
		},
	})
	if err != nil {
		return fmt.Errorf("repo.DeleteUser: query user items: %w", err)
	}
	if len(items) == 0 {
		return nil
	}

	// Extract PK+SK from each item to build the batch delete key list.
	keys := make([]map[string]ddbTypes.AttributeValue, 0, len(items))
	for _, item := range items {
		pk, hasPK := item["PK"]
		sk, hasSK := item["SK"]
		if !hasPK || !hasSK {
			continue
		}
		keys = append(keys, map[string]ddbTypes.AttributeValue{"PK": pk, "SK": sk})
	}

	if err := dynamodb.DeleteItem(ctx, table, nil, true, keys, nil); err != nil {
		return fmt.Errorf("repo.DeleteUser: batch delete: %w", err)
	}
	return nil
}

// GetUserAddresses retrieves all saved addresses for a user.
// TODO: implement using dynamodb.Query with begins_with(SK, "ADDR#") once schema is finalized.
func GetUserAddresses(ctx context.Context, _ string) ([]models.Address, error) {
	return nil, fmt.Errorf("repo.GetUserAddresses: not implemented")
}

// GetUserPreferences retrieves preference settings for a user.
// TODO: implement once schema is finalized.
func GetUserPreferences(ctx context.Context, userID string) (*models.Preferences, error) {
	key, err := dynamodb.BuildKey("PK", "USER#"+userID, "SK", "PREFS")
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserPreferences: build key: %w", err)
	}

	var prefs models.Preferences
	if err := dynamodb.GetItemAs(ctx, table, key, false, nil, &prefs); err != nil {
		return nil, fmt.Errorf("repo.GetUserPreferences: %w", err)
	}
	return &prefs, nil
}
