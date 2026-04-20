package main

import (
	"komodo-auth-api/internal/config"
	"komodo-auth-api/internal/handlers"
	"komodo-auth-api/internal/oauth/clients"
	"net/http"
	"os"
	"time"

	awsEC "github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// Order matters: SM must run before JWT (needs JWT_* keys) and before ElastiCache (needs endpoint).
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
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_AUDIENCE,
			config.JWT_ISSUER,
			config.JWT_KID,
			config.REGISTERED_CLIENTS,
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

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	ecCfg := awsEC.Config{
		Endpoint: os.Getenv(config.AWS_ELASTICACHE_ENDPOINT),
		Password: os.Getenv(config.AWS_ELASTICACHE_PASSWORD),
		DB:       os.Getenv(config.AWS_ELASTICACHE_DB),
	}
	if err := awsEC.Init(ecCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	if err := clients.Load(); err != nil {
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
	mux.HandleFunc("GET /health", health.HealthHandler)
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

	srv.Run(server, os.Getenv(config.PORT), 30*time.Second)
}
