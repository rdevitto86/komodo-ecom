package main

// TODO: implement returns/RMA service
//
// Deployment target: AWS Lambda + DynamoDB.
// Start in order-api when volume is low; migrate here when return seams become painful.
//
// Startup sequence:
//   init() → logger → secrets manager → dynamodb client
//   main() → lambda.Start(adapter.ProxyWithContext) or http.ListenAndServe for local
//
// Key responsibilities:
//   - Accept return requests against a confirmed order
//   - Enforce return window (configurable per product/category)
//   - Track RMA lifecycle: requested → approved → received → processed
//   - On approval: trigger refund via payments-api
//   - On receipt: trigger restock via inventory-api
//   - Reverse loyalty points via loyalty-api
//   - Notify customer via notifications-api (or events-api fan-out)
//
// Integration points (all internal JWT):
//   - order-api: validate order exists, belongs to user, is within return window
//   - payments-api: POST /refunds
//   - inventory-api: POST /stock/{sku}/restock
//   - loyalty-api: POST /me/points/reverse
//   - notifications-api: return status updates
