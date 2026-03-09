package dynamodb

import (
	"context"
	"fmt"
	logger "komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type QueryInput struct {
	TableName              string
	IndexName              *string
	KeyConditionExpression string
	FilterExpression       *string
	ExpressionValues       map[string]types.AttributeValue
	ExpressionNames        map[string]string
	Limit                  *int32
	ScanIndexForward       *bool
	ExclusiveStartKey      map[string]types.AttributeValue
}

type QueryOutput struct {
	Items            []map[string]types.AttributeValue
	LastEvaluatedKey map[string]types.AttributeValue
	Count            int32
}

type ScanInput struct {
	TableName         string
	IndexName         *string
	FilterExpression  *string
	ExpressionValues  map[string]types.AttributeValue
	ExpressionNames   map[string]string
	Limit             *int32
	ExclusiveStartKey map[string]types.AttributeValue
}

type ScanOutput struct {
	Items            []map[string]types.AttributeValue
	LastEvaluatedKey map[string]types.AttributeValue
	Count            int32
}

// Queries DynamoDB and returns the result.
func Query(ctx context.Context, input QueryInput) (*QueryOutput, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "Query")
	}

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(input.TableName),
		KeyConditionExpression: aws.String(input.KeyConditionExpression),
	}

	// Optional parameters
	if input.IndexName != nil {
		queryInput.IndexName = input.IndexName
	}
	if input.FilterExpression != nil {
		queryInput.FilterExpression = input.FilterExpression
	}
	if input.ExpressionValues != nil {
		queryInput.ExpressionAttributeValues = input.ExpressionValues
	}
	if input.ExpressionNames != nil {
		queryInput.ExpressionAttributeNames = input.ExpressionNames
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.ScanIndexForward != nil {
		queryInput.ScanIndexForward = input.ScanIndexForward
	}
	if input.ExclusiveStartKey != nil {
		queryInput.ExclusiveStartKey = input.ExclusiveStartKey
	}

	// Execute query
	result, err := client.Query(ctx, queryInput)
	if err != nil {
		logger.Error("dynamodb failed to query", err)
		return nil, WrapError(err, "Query")
	}

	return &QueryOutput{
		Items:            result.Items,
		LastEvaluatedKey: result.LastEvaluatedKey,
		Count:            result.Count,
	}, nil
}

// Unmarshals the query result into the provided output interface.
func QueryAs(ctx context.Context, input QueryInput, out interface{}) (*QueryOutput, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "QueryAs")
	}

	// Execute query
	result, err := Query(ctx, input)

	if err != nil {
		logger.Error("dynamodb failed to query", err)
		return nil, WrapError(err, "QueryAs query")
	}
	if err := attributevalue.UnmarshalListOfMaps(result.Items, out); err != nil {
		logger.Error("dynamodb failed to unmarshal items", err)
		return nil, WrapError(err, "QueryAs unmarshal")
	}
	return result, nil
}

// Queries DynamoDB and returns all items.
func QueryAll(ctx context.Context, input QueryInput) ([]map[string]types.AttributeValue, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "QueryAll")
	}

	var allItems []map[string]types.AttributeValue
	var lastKey map[string]types.AttributeValue

	for {
		input.ExclusiveStartKey = lastKey

		// Execute query
		result, err := Query(ctx, input)
		if err != nil {
			logger.Error("dynamodb failed to query", err)
			return nil, WrapError(err, "QueryAll query")
		}

		allItems = append(allItems, result.Items...)

		if result.LastEvaluatedKey == nil { break }
		lastKey = result.LastEvaluatedKey
	}
	return allItems, nil
}

// Unmarshals the query result into the provided output interface.
func QueryAllAs(ctx context.Context, input QueryInput, out interface{}) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "QueryAllAs")
	}

	// Execute query
	items, err := QueryAll(ctx, input)

	if err != nil {
		logger.Error("dynamodb failed to query", err)
		return WrapError(err, "QueryAllAs query")
	}
	if err = attributevalue.UnmarshalListOfMaps(items, out); err != nil {
		logger.Error("dynamodb failed to unmarshal items", err)
		return WrapError(err, "QueryAllAs unmarshal")
	}
	return nil
}

// Scans DynamoDB and returns the result.
func Scan(ctx context.Context, input ScanInput) (*ScanOutput, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "Scan")
	}

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(input.TableName),
	}

	if input.IndexName != nil {
		scanInput.IndexName = input.IndexName
	}
	if input.FilterExpression != nil {
		scanInput.FilterExpression = input.FilterExpression
	}
	if input.ExpressionValues != nil {
		scanInput.ExpressionAttributeValues = input.ExpressionValues
	}
	if input.ExpressionNames != nil {
		scanInput.ExpressionAttributeNames = input.ExpressionNames
	}
	if input.Limit != nil {
		scanInput.Limit = input.Limit
	}
	if input.ExclusiveStartKey != nil {
		scanInput.ExclusiveStartKey = input.ExclusiveStartKey
	}

	// Execute scan
	result, err := client.Scan(ctx, scanInput)
	if err != nil {
		logger.Error("dynamodb failed to scan", err)
		return nil, WrapError(err, "Scan")
	}

	return &ScanOutput{
		Items:            result.Items,
		LastEvaluatedKey: result.LastEvaluatedKey,
		Count:            result.Count,
	}, nil
}

// Unmarshals the scan result into the provided output interface.
func ScanAs(ctx context.Context, input ScanInput, out interface{}) (*ScanOutput, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "ScanAs")
	}

	// Execute scan
	result, err := Scan(ctx, input)

	if err != nil {
		logger.Error("dynamodb failed to scan", err)
		return nil, WrapError(err, "ScanAs scan")
	} 
	if err = attributevalue.UnmarshalListOfMaps(result.Items, out); err != nil {
		logger.Error("dynamodb failed to unmarshal items", err)
		return nil, WrapError(err, "ScanAs unmarshal")
	}
	return result, nil
}

// Scans DynamoDB and returns all items.
func ScanAll(ctx context.Context, input ScanInput) ([]map[string]types.AttributeValue, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "ScanAll")
	}

	var allItems []map[string]types.AttributeValue
	var lastKey map[string]types.AttributeValue

	for {
		input.ExclusiveStartKey = lastKey

		// Execute scan
		result, err := Scan(ctx, input)
		if err != nil {
			logger.Error("dynamodb failed to scan", err)
			return nil, WrapError(err, "ScanAll scan")
		}

		allItems = append(allItems, result.Items...)

		if result.LastEvaluatedKey == nil { break }
		lastKey = result.LastEvaluatedKey
	}
	return allItems, nil
}

// Unmarshals the scan result into the provided output interface.
func ScanAllAs(ctx context.Context, input ScanInput, out interface{}) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "ScanAllAs")
	}

	// Execute scan
	items, err := ScanAll(ctx, input)

	if err != nil {
		logger.Error("dynamodb failed to scan all", err)
		return WrapError(err, "ScanAllAs scan")
	}
	if err = attributevalue.UnmarshalListOfMaps(items, out); err != nil {
		logger.Error("dynamodb failed to unmarshal items", err)
		return WrapError(err, "ScanAllAs unmarshal")
	}
	return nil
}
