package cdc

import (
	komodoEvents "komodo-forge-sdk-go/events"

	"github.com/aws/aws-lambda-go/events"
)

// ClassifyResult holds the classification of a DynamoDB stream record as a
// business event. Domain is the short name used to build the SNS topic ARN
// (e.g. "order" → komodo-order-events-prod.fifo).
type ClassifyResult struct {
	EventType  komodoEvents.EventType
	Source     komodoEvents.Source
	EntityID   string
	EntityType komodoEvents.EntityType
	Domain     string
	Payload    map[string]any
}

// Classifier inspects old/new DynamoDB attribute images for a single stream
// record and returns either a ClassifyResult or (zero, false) if the change is
// not business-meaningful and should be silently skipped.
type Classifier func(
	eventName string,
	old, new map[string]events.DynamoDBAttributeValue,
) (ClassifyResult, bool)

var registry = map[string]Classifier{}

// RegisterClassifier registers a domain classifier for the given DynamoDB
// table name. Typically called from a domain package's init().
func RegisterClassifier(tableName string, fn Classifier) {
	registry[tableName] = fn
}

func classify(
	tableName, eventName string,
	old, new map[string]events.DynamoDBAttributeValue,
) (ClassifyResult, bool) {
	fn, ok := registry[tableName]
	if !ok {
		return ClassifyResult{}, false
	}
	return fn(eventName, old, new)
}
