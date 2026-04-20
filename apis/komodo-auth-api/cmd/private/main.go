package main

import (
	"komodo-auth-api/internal/config"
	"komodo-auth-api/internal/handlers"
	"komodo-auth-api/internal/oauth/clients"
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment.
// Internal needs JWT keys (for ValidateTokenHandler) and the client clients.
// ElastiCache is not needed — revocation checks are public-port concerns.
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
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_AUDIENCE,
			config.JWT_ISSUER,
			config.JWT_KID,
			config.REGISTERED_CLIENTS,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	if err := clients.Load(); err != nil {
		logger.Fatal("failed to load client registry", err)
		os.Exit(1)
	}

	logger.Info("auth-api internal: bootstrap complete")
}

func main() {
	// Internal stack — no JWT auth required. VPC/IAM is the access control layer.
	// Requiring a JWT on /internal/token/validate would create a circular dependency:
	// a service would need a valid token to check if its token is valid.
	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	mux.Handle("POST /internal/token/validate", mw.Chain(http.HandlerFunc(handlers.ValidateTokenHandler), internalMW...))
	mux.Handle("GET /internal/clients", mw.Chain(http.HandlerFunc(handlers.ListClientsHandler), internalMW...))
	mux.Handle("GET /internal/clients/{id}", mw.Chain(http.HandlerFunc(handlers.GetClientHandler), internalMW...))

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
