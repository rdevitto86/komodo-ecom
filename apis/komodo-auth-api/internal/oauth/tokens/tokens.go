// Package tokens manages JWT revocation state in ElastiCache (Redis).
//
// Key scheme:
//
//	issued:jti:<jti>  — written on token issue; TTL = token lifetime (so it
//	                    auto-purges when the token would have expired)
//	revoked:jti:<jti> — written on revocation; TTL = remaining token lifetime
//
// A token is considered revoked if a "revoked:jti:<jti>" key exists in Redis.
// The "issued:jti:<jti>" key is informational only; it lets us prove that a
// JTI was issued by this service before accepting a revoke request for it.
package tokens

import (
	"fmt"
	"time"

	awsEC "github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

const (
	issuedKeyPrefix  = "issued:jti:"
	revokedKeyPrefix = "revoked:jti:"
	sentinelValue    = "1" // sentinel value stored at each key
)

// StoreIssued records a newly-issued JTI so it can be validated against later.
// ttl is the full token lifetime in seconds.
func StoreIssued(jti string, ttl int64) error {
	if jti == "" {
		return fmt.Errorf("tokencache: StoreIssued called with empty JTI")
	}

	key := issuedKeyPrefix + jti
	if err := awsEC.Set(key, sentinelValue, ttl); err != nil {
		return fmt.Errorf("tokencache: failed to store issued JTI %q: %w", jti, err)
	}

	logger.Info("tokencache: stored issued JTI " + jti)
	return nil
}

// StoreRevoked marks a JTI as revoked with TTL = remaining token lifetime.
// ttl must be positive; pass time.Until(expiresAt) converted to seconds.
func StoreRevoked(jti string, remaining time.Duration) error {
	if jti == "" {
		return fmt.Errorf("tokencache: StoreRevoked called with empty JTI")
	}
	if remaining <= 0 { return nil } // Token already expired — nothing to revoke.

	ttlSec := max(int64(remaining.Seconds()), 1)
	key := revokedKeyPrefix + jti
	if err := awsEC.Set(key, sentinelValue, ttlSec); err != nil {
		return fmt.Errorf("tokencache: failed to store revoked JTI %q: %w", jti, err)
	}

	logger.Info("tokencache: stored revoked JTI " + jti)
	return nil
}

// IsRevoked returns true if the JTI has been explicitly revoked.
// Returns (false, nil) on a cache miss — a miss means the token is still valid
// from a revocation perspective (expiry is checked separately by the caller).
func IsRevoked(jti string) (bool, error) {
	if jti == "" { return false, nil } // No JTI means we cannot check revocation; treat as not revoked.

	key := revokedKeyPrefix + jti
	val, err := awsEC.Get(key)
	if err != nil {
		return false, fmt.Errorf("tokencache: failed to check revocation for JTI %q: %w", jti, err)
	}

	return val != "", nil
}
