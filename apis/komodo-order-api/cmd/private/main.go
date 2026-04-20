package main

import (
	"net/http"
	"os"
	"time"

	"komodo-order-api/internal/config"
	"komodo-order-api/internal/handlers"

	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"

	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init mirrors the public server bootstrap so this binary can run independently.
// Secrets, DynamoDB, and JWT keys are shared resources — both binaries need them.
func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)
	logger.Info("order-api private: bootstrap complete")
}

func main() {
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

	// Internal order lookup — for returns and payments-api.
	// TODO: implement GetOrderInternal handler once service layer supports it.

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
