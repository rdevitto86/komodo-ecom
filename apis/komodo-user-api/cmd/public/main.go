package main

import (
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamodb"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/config"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"net/http"
	"os"
	"time"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// AWS client bootstrapping lives here so warm Lambda invocations skip it entirely.
func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)

	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
		Keys: []string{
			"DYNAMODB_ENDPOINT",
			"DYNAMODB_ACCESS_KEY",
			"DYNAMODB_SECRET_KEY",
			"DYNAMODB_TABLE",
			"USER_API_CLIENT_ID",
			"USER_API_CLIENT_SECRET",
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"MAX_CONTENT_LENGTH",
			"IDEMPOTENCY_TTL_SEC",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			"BUCKET_TTL_SECOND",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	ddbCfg := dynamodb.Config{
		Region:    config.GetConfigValue("AWS_REGION"),
		Endpoint:  config.GetConfigValue("DYNAMODB_ENDPOINT"),
		AccessKey: config.GetConfigValue("DYNAMODB_ACCESS_KEY"),
		SecretKey: config.GetConfigValue("DYNAMODB_SECRET_KEY"),
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
	mux.HandleFunc("GET /health", handlers.HealthHandler)

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

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, config.GetConfigValue("PORT"), 30*time.Second)
}
