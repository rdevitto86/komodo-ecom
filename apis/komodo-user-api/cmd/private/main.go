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
// The internal function requests a narrower secret set than public — no rate limiter,
// CSRF, or idempotency keys needed.
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
	mux.HandleFunc("GET /health", health.HealthHandler)

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

	srv.Run(server, os.Getenv(config.PORT_PRIVATE), 30*time.Second)
}
