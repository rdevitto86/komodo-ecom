package main

import (
	"net/http"

	"komodo-loyalty-api/internal/handlers"

	health "github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
)

func main() {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", health.HealthHandler)

	// Reviews — absorbed from komodo-reviews-api
	mux.HandleFunc("POST /me/reviews", handlers.SubmitReview)
	mux.HandleFunc("PUT /me/reviews/{reviewId}", handlers.UpdateReview)
	mux.HandleFunc("DELETE /me/reviews/{reviewId}", handlers.DeleteReview)
	mux.HandleFunc("GET /items/{itemId}/reviews", handlers.ListItemReviews)

	// TODO: wire bootstrap (logger, secrets, DynamoDB, Redis) and middleware stack
	// before enabling loyalty-specific routes (GET /me/loyalty, POST /me/loyalty/redeem).
	_ = mux
}
