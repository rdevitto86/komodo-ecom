package main

import (
	"komodo-auth-api/internal/handlers"
	"komodo-auth-api/internal/registry"
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment.
// Internal needs JWT keys (for ValidateTokenHandler) and the client registry.
// ElastiCache is not needed — revocation checks are public-port concerns.
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
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"REGISTERED_CLIENTS",
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

	if err := registry.Load(); err != nil {
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
	mux.HandleFunc("GET /health", handlers.HealthHandler)

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

	srv.Run(server, os.Getenv("INTERNAL_PORT"), 30*time.Second)
}
