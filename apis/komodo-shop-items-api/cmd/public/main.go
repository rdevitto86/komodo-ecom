package main

import (
	"komodo-shop-items-api/internal/config"
	"komodo-shop-items-api/internal/handlers"
	"net/http"
	"os"
	"time"

	awsS3 "github.com/rdevitto86/komodo-forge-sdk-go/aws/s3"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
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
}

func main() {
	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys: []string{
			config.S3_ENDPOINT,
			config.S3_ACCESS_KEY,
			config.S3_SECRET_KEY,
			config.S3_ITEMS_BUCKET,
			config.SHOP_ITEMS_API_CLIENT_ID,
			config.SHOP_ITEMS_API_CLIENT_SECRET,
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
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	} else {
		logger.Info("aws secrets manager initialized successfully")
	}

	s3Cfg := awsS3.Config{
		Region:    os.Getenv(config.AWS_REGION),
		Endpoint:  os.Getenv(config.S3_ENDPOINT),
		AccessKey: os.Getenv(config.S3_ACCESS_KEY),
		SecretKey: os.Getenv(config.S3_SECRET_KEY),
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
	mux.HandleFunc("GET /health", health.HealthHandler)

	mux.Handle("GET /item/inventory", mw.Chain(http.HandlerFunc(handlers.GetInventory), itemMW...))
	mux.Handle("GET /item/{sku}", mw.Chain(http.HandlerFunc(handlers.GetItemBySKU), itemMW...))
	mux.Handle("POST /item/suggestion", mw.Chain(http.HandlerFunc(handlers.GetSuggestions), protectedMW...))

	mux.Handle("GET /services/repair", mw.Chain(http.HandlerFunc(handlers.GetRepairServices), itemMW...))
	mux.Handle("GET /services/repair/{id}", mw.Chain(http.HandlerFunc(handlers.GetRepairService), itemMW...))

	server := &http.Server{
		Addr:              ":" + os.Getenv(config.PORT),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT), 30*time.Second)
}
