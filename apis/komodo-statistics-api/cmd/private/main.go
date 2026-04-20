package main

import (
	"net/http"
	"os"
	"time"

	"komodo-statistics-api/internal/config"
	"komodo-statistics-api/internal/handlers"

	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init mirrors the public server bootstrap so this binary can run independently.
// Full secrets + DB bootstrap should be added here when the first internal handler
// is implemented (see public/main.go for the complete init pattern).
func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)
	logger.Info("statistics-api private: bootstrap complete")
}

func main() {
	// Private middleware stack: no CORS, CSRF, rate-limiting, or sanitization.
	// Auth enforces JWT validity; RequireServiceScope gates to service-to-service calls.
	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.AuthMiddleware,
		mw.RequireServiceScope,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	// Event ingestion — event-bus-api posts domain events here to drive stat updates.
	mux.Handle("POST /internal/events", mw.Chain(handlers.EventHandler, internalMW...))

	// Admin / inter-service stat reads.
	mux.Handle("GET /internal/stats/dashboard", mw.Chain(handlers.DashboardHandler, internalMW...))
	mux.Handle("GET /internal/stats/items/{itemId}", mw.Chain(handlers.ItemStatsHandler, internalMW...))

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
