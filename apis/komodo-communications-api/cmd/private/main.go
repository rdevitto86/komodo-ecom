package main

// TODO: implement communications-api private server
//
// Private routes (internal JWT, no CORS/CSRF — called by other services only):
//   POST /send/email  — transactional email via provider (SES/SendGrid)
//   POST /send/sms    — SMS via provider (SNS/Twilio)
//   POST /send/push   — in-app push notification
