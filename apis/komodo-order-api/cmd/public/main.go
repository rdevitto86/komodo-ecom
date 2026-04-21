package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"komodo-order-api/internal/config"
	"komodo-order-api/internal/handlers"
	"komodo-order-api/internal/service"

	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamo"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secretsmanager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// AWS client bootstrapping lives here so warm Lambda invocations skip it entirely.
func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)

	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys: []string{
			config.AWS_ELASTICACHE_ENDPOINT,
			config.AWS_ELASTICACHE_PASSWORD,
			config.AWS_ELASTICACHE_DB,
			config.DYNAMODB_ORDERS_TABLE,
			config.DYNAMODB_ACCESS_KEY,
			config.DYNAMODB_SECRET_KEY,
			config.DYNAMODB_ENDPOINT,
			config.CART_API_URL,
			config.INVENTORY_API_URL,
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_ISSUER,
			config.JWT_AUDIENCE,
			config.JWT_KID,
			config.MAX_CONTENT_LENGTH,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
			config.IDEMPOTENCY_TTL_SEC,
			config.IP_WHITELIST,
			config.IP_BLACKLIST,
			config.BUCKET_TTL_SECOND,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	ddbCfg := dynamo.Config{
		Region:    os.Getenv(config.AWS_REGION),
		Endpoint:  os.Getenv(config.DYNAMODB_ENDPOINT),
		AccessKey: os.Getenv(config.DYNAMODB_ACCESS_KEY),
		SecretKey: os.Getenv(config.DYNAMODB_SECRET_KEY),
	}
	if err := dynamo.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	}

	eCfg := elasticache.Config{
		Endpoint: os.Getenv(config.AWS_ELASTICACHE_ENDPOINT),
		Password: os.Getenv(config.AWS_ELASTICACHE_PASSWORD),
		DB:       os.Getenv(config.AWS_ELASTICACHE_DB),
	}
	if err := elasticache.Init(eCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	// order-api validates incoming user JWTs. Both public and private keys are loaded
	// so the service can verify tokens issued by auth-api.
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("order-api public: bootstrap complete")
}

func main() {
	// order-api has no real adapters yet — pass nil so the service skips
	// cart-api token validation and inventory hold confirmation until the
	// HTTP adapter implementations land.
	//
	// TODO: wire real adapters once cart-api and shop-inventory-api HTTP clients
	// are implemented under internal/adapters/.
	// nil adapters: cart, inventory, and user service adapters are not yet wired.
	// TODO: wire real adapters once HTTP clients land under internal/adapters/.
	orderSvc := service.NewOrderService(nil, nil, nil)

	writeMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
		mw.IdempotencyMiddleware,
	}

	readMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
	}

	// guestWriteMW is the same as writeMW but without AuthMiddleware, allowing
	// unauthenticated (guest) callers through. The handler validates identity
	// via optional JWT context or email in the request body.
	guestWriteMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
		mw.IdempotencyMiddleware,
	}

	// guestReadMW is the read stack without AuthMiddleware — JWT is optional.
	// Used for GET /orders/{orderId} which supports both authenticated and guest access.
	guestReadMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	// Unified order submission — optional JWT; email required for guests.
	mux.Handle("POST /orders", mw.Chain(handlers.PlaceOrderUnified(orderSvc), guestWriteMW...))

	// Guest-compatible unified order lookup — no JWT required.
	// Registered before /me/orders/{orderId} to avoid ServeMux ambiguity.
	mux.Handle("GET /orders/{orderId}", mw.Chain(handlers.GetOrderUnified(orderSvc), guestReadMW...))

	// Authenticated order routes — require JWT.
	mux.Handle("POST /me/orders", mw.Chain(handlers.PlaceOrder(orderSvc), writeMW...))
	mux.Handle("GET /me/orders", mw.Chain(handlers.ListOrders(orderSvc), readMW...))
	mux.Handle("GET /me/orders/{orderId}", mw.Chain(handlers.GetOrder(orderSvc), readMW...))
	mux.Handle("POST /me/orders/{orderId}/cancel", mw.Chain(handlers.CancelOrder(orderSvc), writeMW...))

	// Authenticated returns (RMA) routes — registered before /me/orders/{orderId} to
	// prevent the wildcard pattern from consuming the literal "returns" segment.
	mux.Handle("GET /me/orders/returns", mw.Chain(handlers.ListReturns(), readMW...))
	mux.Handle("POST /me/orders/returns", mw.Chain(handlers.CreateReturn(), writeMW...))
	mux.Handle("GET /me/orders/returns/{returnId}", mw.Chain(handlers.GetReturn(), readMW...))
	mux.Handle("DELETE /me/orders/returns/{returnId}", mw.Chain(handlers.CancelReturn(), writeMW...))

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT), 30*time.Second)
}

// mustParseInt64 parses s as int64. Returns fallback on empty or parse failure.
func mustParseInt64(s string, fallback int64) int64 {
	if s == "" {
		return fallback
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
