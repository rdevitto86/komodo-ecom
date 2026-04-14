package main

// TODO: implement promotions service
//
// Deployment target: EC2 / ECS Fargate. Called on every cart load and at
// checkout — too frequent for Lambda cold starts.
//
// Startup sequence:
//
//	init() → logger → secrets manager → dynamodb client
//	main() → http.ListenAndServe(PORT, mux)
//
// Key responsibilities:
//   - Validate promo codes at checkout (guest and authenticated)
//   - Surface applicable automatic promotions for a given cart context
//   - Enforce redemption rules: per-user cap, global cap, date window,
//     minimum order value, eligible SKUs/categories
//   - Promo types: percentage_off, fixed_amount_off, free_shipping, buy_x_get_y
//   - CRUD for internal/admin promo management
//
// Integration points (all internal JWT):
//   - cart-api: passes cart contents to POST /promotions/validate for discount calc
//   - order-api: records applied promo_id + discount_cents on order placement
//   - event-bus-api: publishes promotion.redeemed on successful checkout
