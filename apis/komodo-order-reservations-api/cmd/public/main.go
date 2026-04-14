package main

import (
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-order-reservations-api/internal/handlers"
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

	mux.Handle("GET /slots", mw.Chain(http.HandlerFunc(handlers.GetAvailableSlots), slotMW...))
	mux.Handle("GET /slots/{date}", mw.Chain(http.HandlerFunc(handlers.GetSlotsByDate), slotMW...))

	mux.Handle("POST /bookings", mw.Chain(http.HandlerFunc(handlers.CreateBooking), bookingMW...))
	mux.Handle("GET /bookings/{id}", mw.Chain(http.HandlerFunc(handlers.GetBooking), bookingMW...))
	mux.Handle("PUT /bookings/{id}/cancel", mw.Chain(http.HandlerFunc(handlers.CancelBooking), bookingMW...))
	mux.Handle("PUT /bookings/{id}/confirm", mw.Chain(http.HandlerFunc(handlers.ConfirmBooking), bookingMW...))

	// TODO: add internal route for schedule sync (POST /internal/slots/sync)
	// This is called by the external office scheduling system when technician availability changes.
	// Should be on a separate internal port or protected by internal auth middleware.

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
