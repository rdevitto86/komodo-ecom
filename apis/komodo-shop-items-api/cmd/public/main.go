package main

import (
	"komodo-shop-items-api/internal/handlers"
	"net/http"
	"os"
	"time"

	awsS3 "github.com/rdevitto86/komodo-forge-sdk-go/aws/s3"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
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
		Region:    os.Getenv("AWS_REGION"),
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
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

	mux.Handle("GET /item/inventory", mw.Chain(http.HandlerFunc(handlers.GetInventory), itemMW...))
	mux.Handle("GET /item/{sku}", mw.Chain(http.HandlerFunc(handlers.GetItemBySKU), itemMW...))
	mux.Handle("POST /item/suggestion", mw.Chain(http.HandlerFunc(handlers.GetSuggestions), protectedMW...))

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
