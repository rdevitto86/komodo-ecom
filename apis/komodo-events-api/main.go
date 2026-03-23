package main

// TODO: implement event relay service
//
// Deployment target: deferred — build now, deploy when order.placed fans out to 3+ consumers.
// AWS target: SNS (one topic per event type) + SQS (per-consumer queue with DLQ).
// Local dev: runs as a standard HTTP server on PORT; in-process fan-out for local testing.
//
// Startup sequence:
//   init() → logger → secrets manager → SNS/SQS clients
//   main() → http.ListenAndServe (local) or lambda.Start (Lambda variant, TBD)
//
// Key responsibilities:
//   - Accept inbound events from any producer service (POST /events)
//   - Route to the correct SNS topic by event type
//   - Provide a subscription management API for consumer registration (local dev)
//   - Enforce event envelope schema (see docs/event-schema.md)
//
// Consumers of key events:
//   order.placed      → loyalty-api, notifications-api, inventory-api
//   order.cancelled   → inventory-api (release hold), notifications-api
//   payment.confirmed → order-api, notifications-api
//   payment.failed    → cart-api (clear hold), notifications-api
//   stock.low         → notifications-api (restock alert)
