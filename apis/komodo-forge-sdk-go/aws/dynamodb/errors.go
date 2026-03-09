package dynamodb

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Sentinel errors for DynamoDB operations
var ErrClientNotInitialized = fmt.Errorf("dynamodb: client not initialized")

// Wraps AWS DynamoDB errors with descriptive messages
func WrapError(err error, operation string) error {
	if err == nil { return nil }

	// Check for specific AWS DynamoDB error types and format accordingly
	var (
		conditionalCheckErr *types.ConditionalCheckFailedException
		resourceNotFoundErr *types.ResourceNotFoundException
		provisionedThroughputErr *types.ProvisionedThroughputExceededException
		requestLimitErr *types.RequestLimitExceeded
		transactionConflictErr *types.TransactionConflictException
		transactionCanceledErr *types.TransactionCanceledException
		duplicateItemErr *types.DuplicateItemException
		resourceInUseErr *types.ResourceInUseException
		tableNotFoundErr *types.TableNotFoundException
		internalServerErr *types.InternalServerError
		itemCollectionErr *types.ItemCollectionSizeLimitExceededException
		throttlingErr *types.ThrottlingException
	)

	switch {
		case errors.As(err, &conditionalCheckErr):
			return fmt.Errorf("dynamodb: conditional check failed during %s: %w", operation, err)

		case errors.As(err, &resourceNotFoundErr):
			return fmt.Errorf("dynamodb: resource not found during %s: %w", operation, err)

		case errors.As(err, &provisionedThroughputErr):
			return fmt.Errorf("dynamodb: provisioned throughput exceeded during %s: %w", operation, err)

		case errors.As(err, &requestLimitErr):
			return fmt.Errorf("dynamodb: request limit exceeded during %s: %w", operation, err)

		case errors.As(err, &transactionConflictErr):
			return fmt.Errorf("dynamodb: transaction conflict during %s: %w", operation, err)

		case errors.As(err, &transactionCanceledErr):
			return fmt.Errorf("dynamodb: transaction canceled during %s: %w", operation, err)

		case errors.As(err, &duplicateItemErr):
			return fmt.Errorf("dynamodb: duplicate item during %s: %w", operation, err)

		case errors.As(err, &resourceInUseErr):
			return fmt.Errorf("dynamodb: resource in use during %s: %w", operation, err)

		case errors.As(err, &tableNotFoundErr):
			return fmt.Errorf("dynamodb: table not found during %s: %w", operation, err)

		case errors.As(err, &internalServerErr):
			return fmt.Errorf("dynamodb: internal server error during %s: %w", operation, err)

		case errors.As(err, &itemCollectionErr):
			return fmt.Errorf("dynamodb: item collection size exceeded during %s: %w", operation, err)

		case errors.As(err, &throttlingErr):
			return fmt.Errorf("dynamodb: throttled during %s: %w", operation, err)

		default:
			return fmt.Errorf("dynamodb: %s failed: %w", operation, err)
	}
}
