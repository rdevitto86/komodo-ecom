package main

import (
	"komodo-forge-sdk-go/aws/dynamodb"
	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	"komodo-forge-sdk-go/crypto/jwt"
	mw "komodo-forge-sdk-go/http/middleware"
	srv "komodo-forge-sdk-go/http/server"
	logger "komodo-forge-sdk-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"net/http"
	"os"
	"time"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// The internal function requests a narrower secret set than public — no rate limiter,
// CSRF, or idempotency keys needed.
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

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("user-api internal: bootstrap complete")
}

func main() {
	// Internal stack — network-layer security (VPC/IAM) is the primary control.
	// JWT scope check provides defense-in-depth and a caller audit trail.
	// No CORS, CSRF, sanitization, or rate limiting — those are browser concerns.
	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.AuthMiddleware,
		mw.ScopeMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	mux.Handle("GET /users/{id}", mw.Chain(handlers.GetProfile, internalMW...))
	mux.Handle("GET /users/{id}/addresses", mw.Chain(handlers.GetAddresses, internalMW...))
	mux.Handle("GET /users/{id}/preferences", mw.Chain(handlers.GetPreferences, internalMW...))
	mux.Handle("GET /users/{id}/payments", mw.Chain(handlers.GetPayments, internalMW...))

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, config.GetConfigValue("INTERNAL_PORT"), 30*time.Second)
}
