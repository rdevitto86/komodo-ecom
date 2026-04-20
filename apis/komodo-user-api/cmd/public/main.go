package main

import (
	"komodo-user-api/internal/config"
	"komodo-user-api/internal/handlers"
	"net/http"
	"os"
	"time"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamodb"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
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
			config.DYNAMODB_ENDPOINT,
			config.DYNAMODB_ACCESS_KEY,
			config.DYNAMODB_SECRET_KEY,
			config.DYNAMODB_TABLE,
			config.USER_API_CLIENT_ID,
			config.USER_API_CLIENT_SECRET,
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_AUDIENCE,
			config.JWT_ISSUER,
			config.JWT_KID,
			config.IP_WHITELIST,
			config.IP_BLACKLIST,
			config.MAX_CONTENT_LENGTH,
			config.IDEMPOTENCY_TTL_SEC,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
			config.BUCKET_TTL_SECOND,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	ddbCfg := dynamodb.Config{
		Region:    os.Getenv(config.AWS_REGION),
		Endpoint:  os.Getenv(config.DYNAMODB_ENDPOINT),
		AccessKey: os.Getenv(config.DYNAMODB_ACCESS_KEY),
		SecretKey: os.Getenv(config.DYNAMODB_SECRET_KEY),
	}
	if err := dynamodb.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	}

	// User-api validates tokens signed by auth-api — requires the shared RSA public key.
	// JWT_PRIVATE_KEY is included because InitializeKeys() requires both keys to be present.
	// The private key is not used for signing in this service.
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("user-api public: bootstrap complete")
}

func main() {
	publicReadMW := []func(http.Handler) http.Handler{
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

	publicWriteMW := []func(http.Handler) http.Handler{
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	mux.Handle("GET /me/profile", mw.Chain(handlers.GetProfile, publicReadMW...))
	mux.Handle("POST /me/profile", mw.Chain(handlers.CreateUser, publicWriteMW...))
	mux.Handle("PUT /me/profile", mw.Chain(handlers.UpdateProfile, publicWriteMW...))
	mux.Handle("DELETE /me/profile", mw.Chain(handlers.DeleteProfile, publicWriteMW...))

	mux.Handle("GET /me/addresses", mw.Chain(handlers.GetAddresses, publicReadMW...))
	mux.Handle("POST /me/addresses", mw.Chain(handlers.AddAddress, publicWriteMW...))
	mux.Handle("PUT /me/addresses/{id}", mw.Chain(handlers.UpdateAddress, publicWriteMW...))
	mux.Handle("DELETE /me/addresses/{id}", mw.Chain(handlers.DeleteAddress, publicWriteMW...))

	mux.Handle("GET /me/payments", mw.Chain(handlers.GetPayments, publicReadMW...))
	mux.Handle("PUT /me/payments", mw.Chain(handlers.UpsertPayment, publicWriteMW...))
	mux.Handle("DELETE /me/payments/{id}", mw.Chain(handlers.DeletePayment, publicWriteMW...))

	mux.Handle("GET /me/preferences", mw.Chain(handlers.GetPreferences, publicReadMW...))
	mux.Handle("PUT /me/preferences", mw.Chain(handlers.UpdatePreferences, publicWriteMW...))
	mux.Handle("DELETE /me/preferences", mw.Chain(handlers.DeletePreferences, publicWriteMW...))

	mux.Handle("GET /me/wishlist", mw.Chain(handlers.GetWishlist, publicReadMW...))
	mux.Handle("POST /me/wishlist/items", mw.Chain(handlers.AddWishlistItem, publicWriteMW...))
	mux.Handle("DELETE /me/wishlist/items/{itemId}", mw.Chain(handlers.RemoveWishlistItem, publicWriteMW...))
	mux.Handle("GET /me/wishlist/availability", mw.Chain(handlers.GetWishlistAvailability, publicReadMW...))
	mux.Handle("POST /me/wishlist/move-to-cart", mw.Chain(handlers.MoveWishlistToCart, publicWriteMW...))

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
