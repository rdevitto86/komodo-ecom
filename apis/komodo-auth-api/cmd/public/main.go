package main

import (
	"komodo-auth-api/internal/handlers"
	"komodo-auth-api/internal/registry"
	"net/http"
	"os"
	"time"

	awsEC "github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// Order matters: SM must run before JWT (needs JWT_* keys) and before ElastiCache (needs endpoint).
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
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"REGISTERED_CLIENTS",
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

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	ecCfg := awsEC.Config{
		Endpoint: os.Getenv("AWS_ELASTICACHE_ENDPOINT"),
		Password: os.Getenv("AWS_ELASTICACHE_PASSWORD"),
		DB:       os.Getenv("AWS_ELASTICACHE_DB"),
	}
	if err := awsEC.Init(ecCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	if err := registry.Load(); err != nil {
		logger.Fatal("failed to load client registry", err)
		os.Exit(1)
	}

	logger.Info("auth-api: bootstrap complete")
}

func main() {
	oauthMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.NormalizationMiddleware,
		mw.SanitizationMiddleware,
		mw.RuleValidationMiddleware,
	}

	// introspect + revoke require a valid client token
	protectedMW := append(oauthMW, mw.ClientTypeMiddleware, mw.AuthMiddleware)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /.well-known/jwks.json", handlers.JWKSHandler)

	mux.Handle("POST /oauth/token", mw.Chain(http.HandlerFunc(handlers.OAuthTokenHandler), oauthMW...))
	mux.Handle("GET /oauth/authorize", mw.Chain(http.HandlerFunc(handlers.OAuthAuthorizeHandler), oauthMW...))
	mux.Handle("POST /oauth/introspect", mw.Chain(http.HandlerFunc(handlers.OAuthIntrospectHandler), protectedMW...))
	mux.Handle("POST /oauth/revoke", mw.Chain(http.HandlerFunc(handlers.OAuthRevokeHandler), protectedMW...))

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
