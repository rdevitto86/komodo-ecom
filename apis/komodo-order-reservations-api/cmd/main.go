package main

import (
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/config"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-order-reservations-api/internal/handlers"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

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
			"RESERVATIONS_API_CLIENT_ID",
			"RESERVATIONS_API_CLIENT_SECRET",
			"DYNAMODB_SLOTS_TABLE",
			"DYNAMODB_BOOKINGS_TABLE",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			// TODO: add BOOKING_HOLD_TTL_SECONDS when checkout flow (Option A) is implemented
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	}
	logger.Info("aws secrets manager initialized successfully")

	// TODO: initialize DynamoDB client when forge SDK aws/dynamodb package is available
	// See: komodo-forge-sdk-go/aws/dynamodb (planned)

	// Public slot routes — no auth required, rate limited
	slotMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
	}

	// Protected booking routes — auth required
	bookingMW := append(slotMW,
		mw.AuthMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.SanitizationMiddleware,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	mux.Handle("GET /slots", chain(http.HandlerFunc(handlers.GetAvailableSlots), slotMW...))
	mux.Handle("GET /slots/{date}", chain(http.HandlerFunc(handlers.GetSlotsByDate), slotMW...))

	mux.Handle("POST /bookings", chain(http.HandlerFunc(handlers.CreateBooking), bookingMW...))
	mux.Handle("GET /bookings/{id}", chain(http.HandlerFunc(handlers.GetBooking), bookingMW...))
	mux.Handle("PUT /bookings/{id}/cancel", chain(http.HandlerFunc(handlers.CancelBooking), bookingMW...))
	mux.Handle("PUT /bookings/{id}/confirm", chain(http.HandlerFunc(handlers.ConfirmBooking), bookingMW...))

	// TODO: add internal route for schedule sync (POST /internal/slots/sync)
	// This is called by the external office scheduling system when technician availability changes.
	// Should be on a separate internal port or protected by internal auth middleware.

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
