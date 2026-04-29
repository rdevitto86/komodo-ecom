package main

import (
	"net/http"
	"os"
	"time"

	"komodo-order-api/internal/config"
	"komodo-order-api/internal/handlers"
	"komodo-order-api/internal/service"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/dynamo"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secretsmanager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init bootstraps secrets, DynamoDB, and JWT keys for the private server.
// Mirrors public/main.go bootstrap so this binary can run independently.
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
			config.DYNAMODB_ORDERS_TABLE,
			config.DYNAMODB_ACCESS_KEY,
			config.DYNAMODB_SECRET_KEY,
			config.DYNAMODB_ENDPOINT,
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_ISSUER,
			config.JWT_AUDIENCE,
			config.JWT_KID,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	ddbCfg := dynamo.Config{
		Region:    os.Getenv(config.AWS_REGION),
		Endpoint:  os.Getenv(config.DYNAMODB_ENDPOINT),
		AccessKey: os.Getenv(config.DYNAMODB_ACCESS_KEY),
		SecretKey: os.Getenv(config.DYNAMODB_SECRET_KEY),
	}
	if err := dynamo.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("order-api private: bootstrap complete")
}

func main() {
	// Adapters are nil until cart-api and shop-inventory-api HTTP clients land.
	// Internal routes only read orders — no adapter dependency.
	orderSvc := service.NewOrderService(nil, nil, nil, nil)

	// Private middleware stack: no CORS, CSRF, rate-limiting, or sanitization.
	// Auth enforces JWT validity; RequireServiceScope enforces service-to-service scope claims.
	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.AuthMiddleware,
		mw.RequireServiceScope,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	// Internal order lookup — no user ownership check; scope-checked JWT only.
	mux.Handle("GET /internal/orders/{orderId}", mw.Chain(handlers.GetOrderInternal(orderSvc), internalMW...))

	// Internal returns (RMA) routes — scope-checked JWT only.
	mux.Handle("GET /internal/returns/{returnId}", mw.Chain(handlers.GetReturnInternal(), internalMW...))
	mux.Handle("PUT /internal/returns/{returnId}/approve", mw.Chain(handlers.ApproveReturn(), internalMW...))
	mux.Handle("PUT /internal/returns/{returnId}/receive", mw.Chain(handlers.ReceiveReturn(), internalMW...))
	mux.Handle("PUT /internal/returns/{returnId}/reject", mw.Chain(handlers.RejectReturn(), internalMW...))

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
