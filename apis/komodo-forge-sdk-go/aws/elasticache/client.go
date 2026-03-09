package elasticache

import (
	"context"
	"fmt"
	logger "komodo-forge-sdk-go/logging/runtime"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client   *redis.Client
	initOnce sync.Once
	initErr  error
)

type Config struct {
	Endpoint    string
	Password    string
	DB          string
}

// Initialize Elasticache/Redis client with provided config
func Init(cfg Config) error {
	initOnce.Do(func() {
		logger.Info("initializing elasticache client")

		if cfg.Endpoint == "" {
			logger.Error("elasticache endpoint not provided", fmt.Errorf("elasticache endpoint not provided"))
			initErr = fmt.Errorf("elasticache endpoint not provided")
			return
		}

		// Create redis client
		client = redis.NewClient(&redis.Options{
			Addr: cfg.Endpoint,
			Password: cfg.Password,
			DB: func() int { db, _ := strconv.Atoi(cfg.DB); return db }(),
		})

		// Ping with timeout to verify connectivity
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			logger.Error("failed to ping elasticache", err)
			initErr = err
			return
		}

		logger.Info("elasticache client initialized successfully")
	})
	return initErr
}

// Get retrieves the string value stored at key
func Get(key string) (string, error) {
	if client == nil {
		logger.Error("elasticache client not initialized", fmt.Errorf("elasticache client not initialized"))
		return "", fmt.Errorf("elasticache client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Warn("cache item not found")
		return "", nil
	}
	if err != nil {
		logger.Error("failed to get cache item", err)
		return "", err
	}
	return val, nil
}

// Set stores a value with the provided TTL (in seconds). Use ttl=0 for no expiration
func Set(key string, value string, ttl int64) error {
	if client == nil {
		logger.Error("elasticache client not initialized", fmt.Errorf("elasticache client not initialized"))
		return fmt.Errorf("elasticache client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	var dur time.Duration
	if ttl > 0 {
		dur = time.Duration(ttl) * time.Second
	}

	if err := client.Set(ctx, key, value, dur).Err(); err != nil {
		logger.Error("failed to set cache item", err)
		return err
	}
	return nil
}

// Delete removes a key from the cache
func Delete(key string) error {
	if client == nil {
		logger.Error("elasticache client not initialized", fmt.Errorf("elasticache client not initialized"))
		return fmt.Errorf("elasticache client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	if err := client.Del(ctx, key).Err(); err != nil {
		logger.Error("failed to delete cache item", err)
		return err
	}
	return nil
}

// Close closes the Elasticache client connection
func Close() error {
	if client == nil {
		logger.Warn("elasticache client not initialized - skipping close")
		return nil
	}
	return client.Close()
}

// token bucket Lua script (atomic): returns {allowed, wait_ms}
var tokenBucketScript = redis.NewScript(`
local now = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local burst = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local ttl = tonumber(ARGV[5])

local data = redis.call('HMGET', KEYS[1], 'tokens', 'ts')
local tokens = tonumber(data[1])
local ts = tonumber(data[2])
if tokens == nil then
  tokens = burst
  ts = now
end
local elapsed = (now - ts) / 1000.0
if elapsed < 0 then elapsed = 0 end
local new_tokens = tokens + elapsed * rate
if new_tokens > burst then new_tokens = burst end
local allowed = 0
local wait_ms = 0
if new_tokens >= requested then
  new_tokens = new_tokens - requested
  allowed = 1
else
  local deficit = requested - new_tokens
  if rate > 0 then
	wait_ms = math.ceil((deficit / rate) * 1000)
  else
	wait_ms = 0
  end
end
redis.call('HMSET', KEYS[1], 'tokens', tostring(new_tokens), 'ts', tostring(now))
redis.call('EXPIRE', KEYS[1], ttl)
return {allowed, tostring(wait_ms)}
`)

// AllowDistributed attempts to consume a token from a distributed token bucket
// Returns (allowed, retryAfter, error)
func AllowDistributed(ctx context.Context, key string, rate, burst float64, ttlSec int) (bool, time.Duration, error) {
	if client == nil {
		logger.Error("elasticache client not initialized", fmt.Errorf("elasticache client not initialized"))
		return false, 0, fmt.Errorf("elasticache client not initialized")
	}

	now := time.Now().UnixMilli()
	res, err := tokenBucketScript.Run(ctx, client, []string{key}, now, rate, burst, 1, ttlSec).Result()
	if err != nil {
		logger.Error("failed to execute token bucket script", err)
		return false, 0, err
	}

	// Script returns [allowed, wait_ms]
	arr, ok := res.([]interface{})
	if !ok || len(arr) < 2 {
		logger.Error("unexpected script result", fmt.Errorf("unexpected result: %v", res))
		return false, 0, fmt.Errorf("unexpected script result")
	}

	// Parse allowed (may be number or string)
	var allowed bool
	switch v := arr[0].(type) {
		case int64:
			allowed = v == 1
		case string:
			allowed = v == "1"
		default:
			allowed = false
	}

	// Parse wait time
	var waitMs int64
	switch v := arr[1].(type) {
		case int64:
			waitMs = v
		case string:
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
				waitMs = parsed
			}
	}

	return allowed, time.Duration(waitMs) * time.Millisecond, nil
}
