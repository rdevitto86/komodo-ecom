package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	shopinventory "komodo-cart-api/internal/adapters/shopinventory/v1"
	shopitems "komodo-cart-api/internal/adapters/shopitems/v1"
	"komodo-cart-api/internal/handlers"
	"komodo-cart-api/internal/service"

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
		os.Getenv("APP_NAME"),
		os.Getenv("LOG_LEVEL"),
		os.Getenv("ENV"),
	)

	smCfg := awsSM.Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
		Prefix:   os.Getenv("AWS_SECRET_PREFIX"),
		Batch:    os.Getenv("AWS_SECRET_BATCH"),
		Keys: []string{
			"AWS_ELASTICACHE_ENDPOINT",
			"AWS_ELASTICACHE_PASSWORD",
			"AWS_ELASTICACHE_DB",
			"DYNAMODB_CARTS_TABLE",
			"DYNAMODB_ACCESS_KEY",
			"DYNAMODB_SECRET_KEY",
			"DYNAMODB_ENDPOINT",
			"INVENTORY_API_URL",
			"SHOP_ITEMS_API_URL",
			"CART_GUEST_TTL_SEC",
			"CART_HOLD_TTL_SEC",
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_ISSUER",
			"JWT_AUDIENCE",
			"JWT_KID",
			"MAX_CONTENT_LENGTH",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			"IDEMPOTENCY_TTL_SEC",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"BUCKET_TTL_SECOND",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	ddbCfg := dynamo.Config{
		Region:    os.Getenv("AWS_REGION"),
		Endpoint:  os.Getenv("DYNAMODB_ENDPOINT"),
		AccessKey: os.Getenv("DYNAMODB_ACCESS_KEY"),
		SecretKey: os.Getenv("DYNAMODB_SECRET_KEY"),
	}
	if err := dynamo.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	}

	eCfg := elasticache.Config{
		Endpoint: os.Getenv("AWS_ELASTICACHE_ENDPOINT"),
		Password: os.Getenv("AWS_ELASTICACHE_PASSWORD"),
		DB:       os.Getenv("AWS_ELASTICACHE_DB"),
	}
	if err := elasticache.Init(eCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	// Cart-api validates incoming user JWTs and signs checkout tokens issued to order-api.
	// Both keys are required by InitializeKeys().
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("cart-api public: bootstrap complete")
}

func main() {
	// Wire services.
	guestTTL    := mustParseInt64(os.Getenv("CART_GUEST_TTL_SEC"), 604800)
	checkoutTTL := mustParseInt64(os.Getenv("CART_HOLD_TTL_SEC"), 900)
	shopCli     := shopitems.NewClient(os.Getenv("SHOP_ITEMS_API_URL"))
	invCli      := shopinventory.NewClient(os.Getenv("INVENTORY_API_URL"))
	guestSvc    := service.NewGuestCartService(guestTTL, shopCli, invCli)
	cartSvc     := service.NewCartService(checkoutTTL, shopCli, invCli, guestSvc)

	authReadMW := []func(http.Handler) http.Handler{
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

	authWriteMW := []func(http.Handler) http.Handler{
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

	guestWriteMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
		mw.IdempotencyMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// Authenticated cart routes — require JWT.
	mux.Handle("GET /me/cart", mw.Chain(handlers.GetMyCart(cartSvc), authReadMW...))
	mux.Handle("POST /me/cart/merge", mw.Chain(handlers.MergeGuestCart(cartSvc), authWriteMW...))
	mux.Handle("POST /me/cart/items", mw.Chain(handlers.AddMyCartItem(cartSvc), authWriteMW...))
	mux.Handle("PUT /me/cart/items/{itemId}", mw.Chain(handlers.UpdateMyCartItem(cartSvc), authWriteMW...))
	mux.Handle("DELETE /me/cart/items/{itemId}", mw.Chain(handlers.RemoveMyCartItem(cartSvc), authWriteMW...))
	mux.Handle("DELETE /me/cart", mw.Chain(handlers.ClearMyCart(cartSvc), authWriteMW...))
	mux.Handle("POST /me/cart/checkout", mw.Chain(handlers.InitiateCheckout(cartSvc), authWriteMW...))

	// Guest cart routes — no JWT, session token via X-Session-ID header.
	mux.Handle("POST /cart", mw.Chain(handlers.CreateGuestCart(guestSvc), guestWriteMW...))
	mux.Handle("GET /cart/{cartId}", mw.Chain(handlers.GetGuestCart(guestSvc), guestReadMW...))
	mux.Handle("POST /cart/{cartId}/items", mw.Chain(handlers.AddGuestCartItem(guestSvc), guestWriteMW...))
	mux.Handle("PUT /cart/{cartId}/items/{itemId}", mw.Chain(handlers.UpdateGuestCartItem(guestSvc), guestWriteMW...))
	mux.Handle("DELETE /cart/{cartId}/items/{itemId}", mw.Chain(handlers.RemoveGuestCartItem(guestSvc), guestWriteMW...))
	mux.Handle("DELETE /cart/{cartId}", mw.Chain(handlers.ClearGuestCart(guestSvc), guestWriteMW...))

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv("PORT"), 30*time.Second)
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
