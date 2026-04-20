package handlers

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-loyalty-api/internal/models"
)

// SubmitReview handles POST /me/reviews.
// Requires a verified purchase for the target item.
// TODO: implement — verify purchase via order-api, write review to DynamoDB.
func SubmitReview(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, models.ReviewErr.NotEligible, httpErr.WithDetail("not implemented"))
}

// UpdateReview handles PUT /me/reviews/{reviewId}.
// TODO: implement — validate ownership, update review in DynamoDB.
func UpdateReview(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, models.ReviewErr.NotFound, httpErr.WithDetail("not implemented"))
}

// DeleteReview handles DELETE /me/reviews/{reviewId}.
// TODO: implement — validate ownership, delete review from DynamoDB.
func DeleteReview(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, models.ReviewErr.NotFound, httpErr.WithDetail("not implemented"))
}

// ListItemReviews handles GET /items/{itemId}/reviews.
// Returns paginated reviews and aggregate rating for a product. No authentication required.
// TODO: implement — query DynamoDB, compute avg rating and count.
func ListItemReviews(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotFound, httpErr.WithDetail("not implemented"))
}
