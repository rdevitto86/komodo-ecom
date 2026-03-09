package main

import (
	awsS3 "komodo-forge-sdk-go/aws/s3"
	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	mw "komodo-forge-sdk-go/http/middleware"
	logger "komodo-forge-sdk-go/logging/runtime"
	"komodo-shop-items-api/internal/handlers"
	"net/http"
	"os"
	"time"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

// chain applies middleware in order: first listed = outermost wrapper.
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func main() {
	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
		Keys: []string{
			"S3_ENDPOINT",
			"S3_ACCESS_KEY",
			"S3_SECRET_KEY",
			"S3_ITEMS_BUCKET",
			"SHOP_ITEMS_API_CLIENT_ID",
			"SHOP_ITEMS_API_CLIENT_SECRET",
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
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	} else {
		logger.Info("aws secrets manager initialized successfully")
	}

	s3Cfg := awsS3.Config{
		Region:    config.GetConfigValue("AWS_REGION"),
		Endpoint:  config.GetConfigValue("S3_ENDPOINT"),
		AccessKey: config.GetConfigValue("S3_ACCESS_KEY"),
		SecretKey: config.GetConfigValue("S3_SECRET_KEY"),
	}
	if err := awsS3.Init(s3Cfg); err != nil {
		logger.Fatal("failed to initialize s3", err)
		os.Exit(1)
	} else {
		logger.Info("s3 initialized successfully")
	}

	// Shared middleware stack for /item routes
	itemMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
	}

	// Extended middleware stack for protected /item routes
	protectedMW := append(itemMW,
		mw.AuthMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.SanitizationMiddleware,
		mw.RuleValidationMiddleware,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	mux.Handle("GET /item/inventory", chain(http.HandlerFunc(handlers.GetInventory), itemMW...))
	mux.Handle("GET /item/{sku}", chain(http.HandlerFunc(handlers.GetItemBySKU), itemMW...))

	mux.Handle("POST /item/suggestion", chain(http.HandlerFunc(handlers.GetSuggestions), protectedMW...))

	server := &http.Server{
		Addr:              ":" + config.GetConfigValue("PORT"),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
	logger.Info("server started successfully")
}
