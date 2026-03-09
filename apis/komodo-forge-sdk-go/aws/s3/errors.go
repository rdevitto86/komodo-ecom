package s3

import (
	"fmt"
)

// Sentinel errors for S3 operations
var ErrClientNotInitialized = fmt.Errorf("s3: client not initialized")

// Wraps AWS S3 errors with descriptive messages
func WrapError(err error, operation string) error {
	if err == nil { return nil }
	return fmt.Errorf("s3: %s failed: %w", operation, err)
}
