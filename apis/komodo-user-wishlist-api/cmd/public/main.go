package main

// TODO: implement wishlist service
//
// Deployment target: EC2 / ECS Fargate. Authenticated-only; called when users
// browse and save items. Low write volume, moderate read volume.
//
// Startup sequence:
//
//	init() → logger → secrets manager → dynamodb client
//	main() → http.ListenAndServe(PORT, mux)
//
// Key responsibilities:
//   - Persist per-user wishlist items with no TTL (unlike cart)
//   - Surface current stock availability for wishlist items on demand
//   - Move wishlist items to cart (calls cart-api POST /me/cart/items)
//   - Wishlist items reference shop-items-api item IDs and SKUs
//
// Integration points (all internal JWT):
//   - shop-inventory-api: GET /stock/{sku} for availability checks
//   - cart-api: POST /me/cart/items when moving items to cart
//   - shop-items-api: item metadata validation on add
