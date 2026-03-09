package dynamodb

import (
	"context"
	"fmt"
	"sync"

	logger "komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	client *dynamodb.Client
	once   sync.Once
	mu     sync.RWMutex
	initErr error
)

type Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

// Initialize the DynamoDB client
func Init(config Config) error {
	once.Do(func() {
		if config.Region == "" {
			logger.Error("dynamodb region is required", fmt.Errorf("dynamodb region is required"))
			initErr = fmt.Errorf("dynamodb region is required")
			return
		}

		ctx := context.Background()
		var cfg aws.Config

		if config.AccessKey != "" && config.SecretKey != "" {
			cfg, initErr = awsconfig.LoadDefaultConfig(
				ctx,
				awsconfig.WithRegion(config.Region),
				awsconfig.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(
						config.AccessKey,
						config.SecretKey,
						"",
					),
				),
			)
		} else {
			cfg, initErr = awsconfig.LoadDefaultConfig(
				ctx, awsconfig.WithRegion(config.Region),
			)
		}

		if initErr != nil {
			logger.Error("dynamodb failed to load config", initErr)
			initErr = WrapError(initErr, "Init")
			return
		}

		opts := []func(*dynamodb.Options){}
		if config.Endpoint != "" {
			opts = append(opts, func(dbOpts *dynamodb.Options) {
				dbOpts.BaseEndpoint = aws.String(config.Endpoint)
			})
		}

		mu.Lock()
		client = dynamodb.NewFromConfig(cfg, opts...)
		mu.Unlock()
	})
	return initErr
}

// Check if the DynamoDB client is initialized
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return client != nil
}

const maxBatchSize = 25

// Retrieves a single item or batch of items from DynamoDB
func GetItem(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	batch bool,
	keys []map[string]types.AttributeValue,
) (interface{}, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "GetItem")
	}
	if batch {
		return batchGetItems(ctx, tableName, keys)
	}
	return getItem(ctx, tableName, key)
}

// Retrieves and unmarshals item(s) into the provided output interface
func GetItemAs(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	batch bool,
	keys []map[string]types.AttributeValue,
	out interface{},
) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "GetItemAs")
	}

	// batch flow
	if batch {
		return batchGetItemsAs(ctx, tableName, keys, out)
	}

	// single item flow
	item, err := getItem(ctx, tableName, key)
	if err != nil {
		logger.Error("failed to get item", err)
		return WrapError(err, "GetItemAs get")
	}
	if err = attributevalue.UnmarshalMap(item, out); err != nil {
		logger.Error("failed to unmarshal item", err)
		return WrapError(err, "GetItemAs unmarshal")
	}
	return nil
}

// Writes a single item or batch of items to DynamoDB
func WriteItem(
	ctx context.Context,
	tableName string,
	item map[string]types.AttributeValue,
	batch bool,
	items []map[string]types.AttributeValue,
	condition *string,
) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "WriteItem")
	}
	if batch {
		return batchWriteItems(ctx, tableName, items)
	}
	return putItem(ctx, tableName, item, condition)
}

// Marshals and writes item(s) to DynamoDB
func WriteItemFrom(
	ctx context.Context,
	tableName string,
	item interface{},
	batch bool,
	items interface{},
	condition *string,
) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "WriteItemFrom")
	}

	// batch flow
	if batch {
		av, err := attributevalue.MarshalList(items)
		if err != nil {
			logger.Error("failed to marshal items", err)
			return WrapError(err, "WriteItemFrom marshal batch items")
		}

		avMaps := make([]map[string]types.AttributeValue, len(av))
		for i, it := range av {
			if m, ok := it.(*types.AttributeValueMemberM); ok {
				avMaps[i] = m.Value
			} else {
				return WrapError(fmt.Errorf("item %d is not a map", i), "WriteItemFrom marshal batch items")
			}
		}
		return batchWriteItems(ctx, tableName, avMaps)
	}

	// single item flow
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		logger.Error("failed to marshal item", err)
		return WrapError(err, "WriteItemFrom marshal single item")
	}
	return putItem(ctx, tableName, av, condition)
}

// Updates an item in DynamoDB
func UpdateItem(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	updateExpr string,
	exprValues map[string]types.AttributeValue,
	exprNames map[string]string,
	condition *string,
) (map[string]types.AttributeValue, error) {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return nil, WrapError(ErrClientNotInitialized, "UpdateItem")
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              key,
		UpdateExpression: aws.String(updateExpr),
		ReturnValues:     types.ReturnValueAllNew,
	}

	// optional params
	if exprValues != nil {
		updateInput.ExpressionAttributeValues = exprValues
	}
	if exprNames != nil {
		updateInput.ExpressionAttributeNames = exprNames
	}
	if condition != nil {
		updateInput.ConditionExpression = condition
	}

	// Execute update item
	result, err := client.UpdateItem(ctx, updateInput)
	if err != nil {
		logger.Error("failed to update item", err)
		return nil, WrapError(err, "update item")
	}
	return result.Attributes, nil
}

// Updates and unmarshals the result into the provided output interface
func UpdateItemAs(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	updateExpr string,
	exprValues map[string]types.AttributeValue,
	exprNames map[string]string,
	condition *string,
	out interface{},
) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "UpdateItemAs")
	}

	// Execute update item
	attrs, err := UpdateItem(ctx, tableName, key, updateExpr, exprValues, exprNames, condition)
	if err != nil {
		logger.Error("failed to update item", err)
		return WrapError(err, "UpdateItemAs")
	}
	if err = attributevalue.UnmarshalMap(attrs, out); err != nil {
		logger.Error("failed to unmarshal item", err)
		return WrapError(err, "UpdateItemAs unmarshal")
	}
	return nil
}

// Deletes a single item or batch of items from DynamoDB
func DeleteItem(
	ctx context.Context,
	tableName string,
	key map[string]types.AttributeValue,
	batch bool,
	keys []map[string]types.AttributeValue,
	condition *string,
) error {
	if client == nil {
		logger.Error("dynamodb client not initialized", fmt.Errorf("dynamodb client not initialized"))
		return WrapError(ErrClientNotInitialized, "DeleteItem")
	}
	if batch {
		return batchDeleteItems(ctx, tableName, keys)
	}
	return deleteItem(ctx, tableName, key, condition)
}

// Creates a DynamoDB key from partition and optional sort key
func BuildKey(
	partitionKey string,
	partitionValue interface{},
	sortKey string,
	sortValue interface{},
) (map[string]types.AttributeValue, error) {
	key := make(map[string]types.AttributeValue)

	pkAttr, err := attributevalue.Marshal(partitionValue)
	if err != nil {
		logger.Error("failed to marshal partition key", err)
		return nil, WrapError(err, "BuildKey marshal partition key")
	}

	key[partitionKey] = pkAttr

	if sortKey != "" && sortValue != nil {
		skAttr, err := attributevalue.Marshal(sortValue)
		if err != nil {
			logger.Error("failed to marshal sort key", err)
			return nil, WrapError(err, "BuildKey marshal sort key")
		}
		key[sortKey] = skAttr
	}
	return key, nil
}
