package main

// TODO: implement order-api public server
//
// Public routes (JWT auth, CORS, rate-limiting):
//   POST   /me/orders                   — place order (consumes cart checkout token)
//   GET    /me/orders                   — list authenticated user's orders
//   GET    /me/orders/{orderId}         — order detail
//   POST   /me/orders/{orderId}/cancel  — cancel order
