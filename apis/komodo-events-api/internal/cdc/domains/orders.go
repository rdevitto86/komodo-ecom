package domains

import (
	"komodo-events-api/internal/cdc"
	komodoEvents "github.com/rdevitto86/komodo-forge-sdk-go/events"

	"github.com/aws/aws-lambda-go/events"
)

// ordersTable must match the DynamoDB table name for the orders domain.
// TODO: confirm this matches the actual table name defined in infra/deploy/cfn/infra.yaml
// and infra/local/localstack/init/ once the order-api schema is finalised.
const ordersTable = "komodo-orders"

func init() {
	cdc.RegisterClassifier(ordersTable, classifyOrder)
}

// TODO: add classifiers for payments, users, inventory, and cart domains here
// as those tables are instrumented with DynamoDB Streams. Each follows the
// same pattern: one file per domain, one init() registration per table.

func classifyOrder(
	eventName string,
	old, new map[string]events.DynamoDBAttributeValue,
) (cdc.ClassifyResult, bool) {
	switch eventName {
		case "INSERT":
			// TODO: "pk" is a placeholder — replace with the actual DynamoDB partition
			// key attribute name from apis/komodo-order-api/docs/data-model.md.
			orderID := attrString(new, "pk")
			if orderID == "" { return cdc.ClassifyResult{}, false }

			return cdc.ClassifyResult{
				EventType:  komodoEvents.EventOrderCreated,
				Source:     komodoEvents.SourceOrderAPI,
				EntityID:   orderID,
				EntityType: komodoEvents.EntityOrder,
				Domain:     "order",
				Payload: map[string]any{
					"order_id": orderID,
					"user_id":  attrString(new, "user_id"),
					"status":   attrString(new, "status"),
				},
			}, true

	case "MODIFY":
		// TODO: expand Payload with additional fields (total_cents, item_count,
		// shipping_address_id) as the order-api data model stabilises.
		oldStatus := attrString(old, "status")
		newStatus := attrString(new, "status")
		if oldStatus == newStatus || newStatus == "" {
			return cdc.ClassifyResult{}, false // not a status transition — skip
		}
		orderID := attrString(new, "pk")
		result := cdc.ClassifyResult{
			Source:     komodoEvents.SourceOrderAPI,
			EntityID:   orderID,
			EntityType: komodoEvents.EntityOrder,
			Domain:     "order",
			Payload: map[string]any{
				"order_id":   orderID,
				"user_id":    attrString(new, "user_id"),
				"old_status": oldStatus,
				"new_status": newStatus,
			},
		}

		switch newStatus {
			case "cancelled":
				result.EventType = komodoEvents.EventOrderCancelled
			case "fulfilled":
				result.EventType = komodoEvents.EventOrderFulfilled
			default:
				result.EventType = komodoEvents.EventOrderStatusUpdated
		}
		return result, true
	}

	return cdc.ClassifyResult{}, false
}

// attrString safely reads a String-typed DynamoDB attribute value.
// Returns "" if the key is absent or the attribute is not a String type.
func attrString(image map[string]events.DynamoDBAttributeValue, key string) string {
	v, ok := image[key]
	if !ok || v.DataType() != events.DataTypeString { return "" }
	return v.String()
}
