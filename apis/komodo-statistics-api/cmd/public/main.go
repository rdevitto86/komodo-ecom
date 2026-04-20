package main

import (
	"net/http"
	"os"
	"time"

	"komodo-statistics-api/internal/config"
	"komodo-statistics-api/internal/handlers"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secretsmanager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

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
			config.SQLITE_DB_PATH,
			config.JWT_PUBLIC_KEY,
			config.JWT_ISSUER,
			config.JWT_AUDIENCE,
			config.JWT_KID,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
			config.IP_WHITELIST,
			config.IP_BLACKLIST,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	// Public stat routes are unauthenticated (banner counters), but the JWT key is
	// loaded so the middleware stack can validate tokens on routes that opt into auth.
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("statistics-api public: bootstrap complete")
}

func main() {
	// statsMW is used for public stat reads — anonymous counters served to the UI.
	// No auth or CSRF required; rate-limiting guards against scraping.
	statsMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.RuleValidationMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	// Item-level social proof counters — served as UI banners (no login required).
	mux.Handle("GET /stats/items/{itemId}/in-cart", mw.Chain(handlers.InCartHandler, statsMW...))
	mux.Handle("GET /stats/items/{itemId}/recently-bought", mw.Chain(handlers.RecentlyBoughtHandler, statsMW...))
	mux.Handle("GET /stats/items/{itemId}/frequently-bought-with", mw.Chain(handlers.FrequentlyBoughtWithHandler, statsMW...))
	mux.Handle("GET /stats/trending", mw.Chain(handlers.TrendingHandler, statsMW...))

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
