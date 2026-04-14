package main

// TODO: implement reviews-api public server
//
// Public routes (JWT auth, CORS, rate-limiting):
//   POST   /me/reviews                  — submit a review (verified purchase required)
//   PUT    /me/reviews/{reviewId}       — update own review
//   DELETE /me/reviews/{reviewId}       — delete own review
//   GET    /items/{itemId}/reviews      — paginated review listing (unauthenticated)
