package rateLimiter

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"komodo-forge-sdk-go/aws/elasticache"
	"komodo-forge-sdk-go/config"
)

type bucket struct {
	mu      sync.Mutex
	tokens  float64
	last    time.Time
	created time.Time
}

type Service interface {
	Allow(ctx context.Context, key string) (allowed bool, wait time.Duration, err error)
	GetUsage(ctx context.Context, key string) (used int, remaining int, reset time.Time, err error)
	Reset(ctx context.Context, key string) error
	LoadConfig(cfg Config) error
	ShouldFailOpen() bool
}

type Config struct {
	RPS             float64
	Burst           float64
	BucketTTLSecond int
	FailOpen        *bool
}

var (
	rlOnce    sync.Once
	rps       float64
	burst     float64
	buckets   sync.Map
	evictOnce sync.Once
)

// Allow attempts to consume a token for the given client key
func Allow(ctx context.Context, key string) (allowed bool, wait time.Duration, err error) {
	env := strings.ToLower(config.GetConfigValue("ENV"))

	if env == "prod" || env == "staging" {
		rpsVal, burstVal := rateConfig()
		ttl := config.GetConfigValue("BUCKET_TTL_SECOND")
		ttlSec, _ := strconv.Atoi(ttl)
		return elasticache.AllowDistributed(ctx, key, rpsVal, burstVal, ttlSec)
	}

	// local process bucket
	b := getBucket(key)
	if !b.allow() {
		return false, b.retryAfter(), nil
	}
	return true, 0, nil
}

// Returns simple usage metrics for the given key
func GetUsage(ctx context.Context, key string) (used int, remaining int, reset time.Time, err error) {
	b := getBucket(key)
	// snapshot under lock
	b.mu.Lock()
	tokens := b.tokens
	b.mu.Unlock()

	_, burstVal := rateConfig()
	remaining = int(tokens)
	usedF := burstVal - tokens
	if usedF < 0 {
		usedF = 0
	}
	used = int(usedF)
	// estimate reset as when a token will next be available
	reset = time.Now().Add(b.retryAfter())
	return used, remaining, reset, nil
}

// Reset removes any in-process bucket state for the given key. This does not
// affect any distributed (Elasticache) state.
func Reset(ctx context.Context, key string) error {
	buckets.Delete(key)
	return nil
}

// Programmatically overrides rate limiter settings (RPS/Burst).
func LoadConfig(cfg Config) error {
	if cfg.RPS > 0 { rps = cfg.RPS }
	if cfg.Burst > 0 { burst = cfg.Burst }
	return nil
}

// Checks and updates the bucket token count
func (bkt *bucket) allow() bool {
	rps, burst := rateConfig()
	now := time.Now()
	bkt.mu.Lock()

	// Refill tokens based on elapsed time
	if !bkt.last.IsZero() {
		elapsed := now.Sub(bkt.last).Seconds()
		if elapsed > 0 {
			bkt.tokens += elapsed * rps
			if bkt.tokens > burst {
				bkt.tokens = burst
			}
		}
	} else {
		bkt.tokens = burst
	}

	allowed := false
	if bkt.tokens >= 1 {
		bkt.tokens -= 1
		allowed = true
	}

	bkt.last = now
	bkt.mu.Unlock()

	return allowed
}

// Estimates how long until the next token is available
func (bkt *bucket) retryAfter() time.Duration {
	rps, _ := rateConfig()
	if rps <= 0 { return time.Second }

	bkt.mu.Lock()
	defer bkt.mu.Unlock()
	deficit := 1 - bkt.tokens
	if deficit <= 0 { return 0 }
	secs := deficit / rps
	return time.Duration(secs * float64(time.Second))
}

// Reads and caches rate limit settings from env vars
func rateConfig() (float64, float64) {
	rlOnce.Do(func() {
		parseFloatEnv := func(key string, dflt float64) float64 {
			if val := strings.TrimSpace(config.GetConfigValue(key)); val != "" {
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					return f
				}
			}
			return dflt
		}

		rps = parseFloatEnv("RATE_LIMIT_RPS", 10)    // default 10 req/sec
		burst = parseFloatEnv("RATE_LIMIT_BURST", 20) // default burst 20

		// stricter validation: treat non-positive rps as invalid and reset
		if rps <= 0 { rps = 10 }
		if burst < 1 { burst = 20 }
	})
	return rps, burst
}

// Retrieves or creates a rate limit bucket for the given key
func getBucket(key string) *bucket {
	// ensure the background evictor is running
	evictOnce.Do(startBucketEvictor)

	if v, ok := buckets.Load(key); ok {
		return v.(*bucket)
	}
	bkt := &bucket{tokens: 0, last: time.Time{}, created: time.Now()}
	actual, _ := buckets.LoadOrStore(key, bkt)
	return actual.(*bucket)
}

// Removes idle buckets after configured TTL
func startBucketEvictor() {
	ttlSec := 300
	if val := strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_BUCKET_TTL_SEC")); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i > 0 {
			ttlSec = i
		}
	}

	ttl := time.Duration(ttlSec) * time.Second
	ticker := time.NewTicker(time.Minute)

	// background goroutine to evict old buckets
	go func() {
		for range ticker.C {
			now := time.Now()
			buckets.Range(func(key, val any) bool {
				bucket := val.(*bucket)

				bucket.mu.Lock()
				lastActive := bucket.last
				if lastActive.IsZero() {
					lastActive = bucket.created
				}
				bucket.mu.Unlock()

				if now.Sub(lastActive) > ttl {
					buckets.Delete(key)
				}
				return true
			})
		}
	}()
}

// Decides fail-open vs fail-closed when the distributed store is unavailable
func ShouldFailOpen() bool {
	v := strings.ToLower(strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_FAIL_OPEN")))
	if v == "" { return true }
	return v == "true" || v == "1" || v == "yes"
}
