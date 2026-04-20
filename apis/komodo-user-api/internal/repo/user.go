package repo

import (
	"context"
	"fmt"
	"os"
	"strings"

	"komodo-user-api/internal/config"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"komodo-user-api/internal/models"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamodb"
)

// table is resolved at startup from config.
// Set DYNAMODB_TABLE=komodo-users for LocalStack and production alike —
// the forge DynamoDB client switches endpoints transparently via DYNAMODB_ENDPOINT.
var table = os.Getenv(config.DYNAMODB_TABLE)

// Key schema: PK=USER#<id>, SK=PROFILE|ADDR#<address_id>|PAY#<payment_id>|PREFS
// Full design is in docs/data-model.md.

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

// addressRecord is the internal DynamoDB representation of an address item.
// It carries PK/SK in addition to the public Address fields so that marshalling
// via WriteItemFrom correctly populates the table keys.
type addressRecord struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	models.Address
}

// addrPK returns the partition-key value for a user's address items.
func addrPK(userID string) string { return "USER#" + userID }

// addrSK returns the sort-key value for a specific address.
func addrSK(addressID string) string { return "ADDR#" + addressID }

// CreateAddress writes a new address item under the user's partition.
// If addr.AddressID is empty a new ID is generated in the form "addr_<12hex>"
// and written back onto addr so the caller sees the assigned ID.
// is_default enforcement is the caller's responsibility — the repo writes as-is.
func CreateAddress(ctx context.Context, userID string, addr *models.Address) error {
	if addr.AddressID == "" {
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		addr.AddressID = "addr_" + raw[:12]
	}

	record := addressRecord{
		PK:      addrPK(userID),
		SK:      addrSK(addr.AddressID),
		Address: *addr,
	}
	if err := dynamodb.WriteItemFrom(ctx, table, record, false, nil, nil); err != nil {
		return fmt.Errorf("repo.CreateAddress: %w", err)
	}
	return nil
}

// GetAddress retrieves a single address by its ID.
func GetAddress(ctx context.Context, userID, addressID string) (*models.Address, error) {
	key, err := dynamodb.BuildKey("PK", addrPK(userID), "SK", addrSK(addressID))
	if err != nil {
		return nil, fmt.Errorf("repo.GetAddress: build key: %w", err)
	}

	var record addressRecord
	if err := dynamodb.GetItemAs(ctx, table, key, false, nil, &record); err != nil {
		return nil, fmt.Errorf("repo.GetAddress: %w", err)
	}
	addr := record.Address
	return &addr, nil
}

// GetUserAddresses retrieves all saved addresses for a user.
// Uses a Query with PK=USER#<userID> and begins_with(SK, "ADDR#").
func GetUserAddresses(ctx context.Context, userID string) ([]models.Address, error) {
	var records []addressRecord
	if err := dynamodb.QueryAllAs(ctx, dynamodb.QueryInput{
		TableName:              table,
		KeyConditionExpression: "PK = :pk AND begins_with(SK, :skPrefix)",
		ExpressionValues: map[string]ddbTypes.AttributeValue{
			":pk":       &ddbTypes.AttributeValueMemberS{Value: addrPK(userID)},
			":skPrefix": &ddbTypes.AttributeValueMemberS{Value: "ADDR#"},
		},
	}, &records); err != nil {
		return nil, fmt.Errorf("repo.GetUserAddresses: %w", err)
	}

	addrs := make([]models.Address, len(records))
	for i, r := range records {
		addrs[i] = r.Address
	}
	return addrs, nil
}

// UpdateAddress replaces an address item in full (PutItem semantics).
// The handler always sends the complete address object, so a full replace is
// simpler and avoids complex UpdateItem expressions. See data-model.md for rationale.
// is_default enforcement is the caller's responsibility.
func UpdateAddress(ctx context.Context, userID string, addr models.Address) error {
	if addr.AddressID == "" {
		return fmt.Errorf("repo.UpdateAddress: address_id is required")
	}

	record := addressRecord{
		PK:      addrPK(userID),
		SK:      addrSK(addr.AddressID),
		Address: addr,
	}
	if err := dynamodb.WriteItemFrom(ctx, table, record, false, nil, nil); err != nil {
		return fmt.Errorf("repo.UpdateAddress: %w", err)
	}
	return nil
}

// DeleteAddress removes a single address item.
func DeleteAddress(ctx context.Context, userID, addressID string) error {
	key, err := dynamodb.BuildKey("PK", addrPK(userID), "SK", addrSK(addressID))
	if err != nil {
		return fmt.Errorf("repo.DeleteAddress: build key: %w", err)
	}
	if err := dynamodb.DeleteItem(ctx, table, key, false, nil, nil); err != nil {
		return fmt.Errorf("repo.DeleteAddress: %w", err)
	}
	return nil
}

// paymentRecord is the internal DynamoDB representation of a payment method item.
// It carries PK/SK in addition to the public PaymentMethod fields so that marshalling
// via WriteItemFrom correctly populates the table keys.
// Token is included here so it is persisted to DynamoDB; it is explicitly zeroed out
// on every read path before the PaymentMethod is returned to callers.
type paymentRecord struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`
	models.PaymentMethod
}

// payPK returns the partition-key value for a user's payment items.
func payPK(userID string) string { return "USER#" + userID }

// paySK returns the sort-key value for a specific payment method.
func paySK(paymentID string) string { return "PAY#" + paymentID }

// UpsertPayment writes a payment method item under the user's partition (PutItem semantics).
// If method.PaymentID is empty a new ID is generated in the form "pay_<12hex>"
// and written back onto method so the caller sees the assigned ID.
// is_default enforcement is the caller's responsibility — the repo writes as-is.
func UpsertPayment(ctx context.Context, userID string, method *models.PaymentMethod) error {
	if method.PaymentID == "" {
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		method.PaymentID = "pay_" + raw[:12]
	}

	record := paymentRecord{
		PK:            payPK(userID),
		SK:            paySK(method.PaymentID),
		PaymentMethod: *method,
	}
	if err := dynamodb.WriteItemFrom(ctx, table, record, false, nil, nil); err != nil {
		return fmt.Errorf("repo.UpsertPayment: %w", err)
	}
	return nil
}

// GetPayment retrieves a single payment method by its ID.
// Token is zeroed before returning — it is stored in DynamoDB for internal use
// by the payments-api but must never be exposed through the user-api response path.
func GetPayment(ctx context.Context, userID, paymentID string) (*models.PaymentMethod, error) {
	key, err := dynamodb.BuildKey("PK", payPK(userID), "SK", paySK(paymentID))
	if err != nil {
		return nil, fmt.Errorf("repo.GetPayment: build key: %w", err)
	}

	var record paymentRecord
	if err := dynamodb.GetItemAs(ctx, table, key, false, nil, &record); err != nil {
		return nil, fmt.Errorf("repo.GetPayment: %w", err)
	}

	pm := record.PaymentMethod
	pm.Token = "" // defense-in-depth: never return the processor token via this API
	return &pm, nil
}

// ListPayments retrieves all saved payment methods for a user.
// Uses a Query with PK=USER#<userID> and begins_with(SK, "PAY#").
// Token is zeroed on every record before returning — see GetPayment for rationale.
func ListPayments(ctx context.Context, userID string) ([]models.PaymentMethod, error) {
	var records []paymentRecord
	if err := dynamodb.QueryAllAs(ctx, dynamodb.QueryInput{
		TableName:              table,
		KeyConditionExpression: "PK = :pk AND begins_with(SK, :skPrefix)",
		ExpressionValues: map[string]ddbTypes.AttributeValue{
			":pk":       &ddbTypes.AttributeValueMemberS{Value: payPK(userID)},
			":skPrefix": &ddbTypes.AttributeValueMemberS{Value: "PAY#"},
		},
	}, &records); err != nil {
		return nil, fmt.Errorf("repo.ListPayments: %w", err)
	}

	methods := make([]models.PaymentMethod, len(records))
	for i, r := range records {
		pm := r.PaymentMethod
		pm.Token = "" // defense-in-depth: never return the processor token via this API
		methods[i] = pm
	}
	return methods, nil
}

// DeletePayment removes a single payment method item.
func DeletePayment(ctx context.Context, userID, paymentID string) error {
	key, err := dynamodb.BuildKey("PK", payPK(userID), "SK", paySK(paymentID))
	if err != nil {
		return fmt.Errorf("repo.DeletePayment: build key: %w", err)
	}
	if err := dynamodb.DeleteItem(ctx, table, key, false, nil, nil); err != nil {
		return fmt.Errorf("repo.DeletePayment: %w", err)
	}
	return nil
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
