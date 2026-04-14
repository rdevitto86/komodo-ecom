package main

import (
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-search-api/internal/handlers"
)

func init() {
	logger.Init(os.Getenv("APP_NAME"), os.Getenv("LOG_LEVEL"), os.Getenv("ENV"))
}

func main() {
	smCfg := awsSM.Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
		Prefix:   os.Getenv("AWS_SECRET_PREFIX"),
		Batch:    os.Getenv("AWS_SECRET_BATCH"),
		Keys: []string{
			"SEARCH_API_CLIENT_ID",
			"SEARCH_API_CLIENT_SECRET",
			"TYPESENSE_HOST",
			"TYPESENSE_PORT",
			"TYPESENSE_API_KEY",
			"TYPESENSE_COLLECTION",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	}
	logger.Info("aws secrets manager initialized successfully")

	// TODO(typesense): initialize Typesense client after secrets are loaded.
	// Add dependency: github.com/typesense/typesense-go
	// Client config: host, port, api_key from secrets above.
	// Call repository.InitTypesense(host, port, apiKey, collection) here.
	// Verify collection exists on startup — log warning if not, don't fatal
	// (search will return IndexUnavailable errors until collection is ready).

	// TODO(subscriber): start events-api subscriber in a background goroutine.
	// subscriber.StartShopItemsSubscriber(ctx) listens for shop-item create/update/delete
	// events and syncs them to the Typesense index.
	// Only start after Typesense client is initialized.

	searchMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.Handle("GET /search", mw.Chain(http.HandlerFunc(handlers.Search), searchMW...))

	// TODO(typesense): add index management routes (internal only):
	//   POST /internal/index/sync  — full re-index from shop-items-api (manual trigger)
	//   DELETE /internal/index     — drop and recreate collection (for schema changes)

	server := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv("PORT"), 30*time.Second)
}
