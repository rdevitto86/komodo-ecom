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

func getItem(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
) (map[string]types.AttributeValue, error) {
	// Execute get item
	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: key,
	})

	if err != nil {
		logger.Error("failed to get item", err)
		return nil, WrapError(err, "getItem")
	}
	if result.Item == nil {
		logger.Error("item not found", fmt.Errorf("dynamodb item not found"))
		return nil, WrapError(fmt.Errorf("item not found"), "getItem")
	}
	return result.Item, nil
}

func putItem(
	ctx context.Context,
	tableName string,
	item map[string]types.AttributeValue,
	condition *string,
) error {
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	if condition != nil {
		putInput.ConditionExpression = condition
	}
	// Execute put item
	if _, err := client.PutItem(ctx, putInput); err != nil {
		logger.Error("failed to put item", err)
		return WrapError(err, "putItem")
	}
	return nil
}

func deleteItem(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	condition *string,
) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: key,
	}
	if condition != nil {
		deleteInput.ConditionExpression = condition
	}

	// Execute delete
	if _, err := client.DeleteItem(ctx, deleteInput); err != nil {
		logger.Error("failed to delete item", err)
		return WrapError(err, "deleteItem")
	}
	return nil
}

func batchGetItems(
	ctx context.Context,
	tableName string,
	keys []map[string]types.AttributeValue,
) ([]map[string]types.AttributeValue, error) {
	if len(keys) == 0 {
		logger.Warn("No keys to batch get")
		return nil, nil
	}

	var allItems []map[string]types.AttributeValue
	for i := 0; i < len(keys); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(keys) { end = len(keys) }

		// Execute batch get
		result, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: { Keys: keys[i:end] },
			},
		})

		if err != nil {
			logger.Error("failed to batch get items", err)
			return nil, WrapError(err, "batchGetItems")
		}
		if items, ok := result.Responses[tableName]; ok {
			allItems = append(allItems, items...)
		}
		if len(result.UnprocessedKeys) > 0 {
			logger.Error("batch get has unprocessed keys", fmt.Errorf("batch get has unprocessed keys"))
			return nil, WrapError(fmt.Errorf("batch get has unprocessed keys"), "batchGetItems")
		}
	}
	return allItems, nil
}

func batchGetItemsAs(
	ctx context.Context,
	tableName string,
	keys []map[string]types.AttributeValue,
	out interface{},
) error {
	// Execute batch get
	items, err := batchGetItems(ctx, tableName, keys)

	if err != nil {
		logger.Error("failed to batch get items", err)
		return err
	}
	if err = attributevalue.UnmarshalListOfMaps(items, out); err != nil {
		logger.Error("failed to unmarshal items", err)
		return WrapError(fmt.Errorf("failed to unmarshal items"), "batchGetItemsAs")
	}
	return nil
}

func batchWriteItems(
	ctx context.Context,
	tableName string,
	items []map[string]types.AttributeValue,
) error {
	if len(items) == 0 {
		logger.Warn("No items to batch write")
		return nil
	}

	// Loop through items in batches
	for i := 0; i < len(items); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(items) { end = len(items) }

		// Create write requests
		writeRequests := make([]types.WriteRequest, end-i)
		for j, item := range items[i:end] {
			writeRequests[j] = types.WriteRequest{
				PutRequest: &types.PutRequest{Item: item},
			}
		}

		// Execute batch write
		result, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})
		
		if err != nil {
			logger.Error("failed to batch write items", err)
			return WrapError(err, "batchWriteItem")
		}
		if len(result.UnprocessedItems) > 0 {
			logger.Error("batch write has unprocessed items", fmt.Errorf("batch write has unprocessed items"))
			return WrapError(fmt.Errorf("batch write has unprocessed items"), "batchWriteItem")
		}
	}
	return nil
}

func batchDeleteItems(
	ctx context.Context,
	tableName string,
	keys []map[string]types.AttributeValue,
) error {
	if len(keys) == 0 {
		logger.Warn("No keys to batch delete")
		return nil
	}

	// Loop through keys in batches
	for i := 0; i < len(keys); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(keys) { end = len(keys) }

		// Create write requests
		writeRequests := make([]types.WriteRequest, end-i)
		for j, key := range keys[i:end] {
			writeRequests[j] = types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{Key: key},
			}
		}

		// Execute batch delete
		result, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})

		if err != nil {
			logger.Error("failed to batch delete items", err)
			return WrapError(err, "batchDeleteItem")
		}
		if len(result.UnprocessedItems) > 0 {
			logger.Error("batch delete has unprocessed items", fmt.Errorf("batch delete has unprocessed items"))
			return WrapError(fmt.Errorf("batch delete has unprocessed items"), "batchDeleteItem")
		}
	}
	return nil
}
