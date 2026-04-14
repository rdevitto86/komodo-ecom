package main

import (
	"net/http"
	"os"
	"time"

	"komodo-address-api/internal/handlers"
	"komodo-address-api/internal/provider"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// Order matters: logger first, then SM (loads JWT_* keys), then JWT init.
func init() {
	logger.Init(os.Getenv("APP_NAME"), os.Getenv("LOG_LEVEL"), os.Getenv("ENV"))

	if err := awsSM.Bootstrap(awsSM.Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
		Prefix:   os.Getenv("AWS_SECRET_PREFIX"),
		Batch:    os.Getenv("AWS_SECRET_BATCH"),
		Keys: []string{
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"ADDRESS_PROVIDER_API_KEY",
			"MAX_CONTENT_LENGTH",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
		},
	}); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("address-api: bootstrap complete")
}

func main() {
	stack := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
	}

	providerClient := provider.NewClient(os.Getenv("ADDRESS_PROVIDER_API_KEY"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)
	mux.Handle("POST /addresses/validate", mw.Chain(handlers.Validate(providerClient), stack...))
	mux.Handle("POST /addresses/normalize", mw.Chain(handlers.Normalize(providerClient), stack...))
	mux.Handle("POST /addresses/geocode", mw.Chain(handlers.Geocode(providerClient), stack...))

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
