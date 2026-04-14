package main

// TODO: implement entitlements-api private server
//
// Private routes (scope-checked JWT, no CORS/CSRF — called by order-api):
//   POST   /entitlements      — grant entitlement on purchase
//   DELETE /entitlements/{id} — revoke entitlement on return/cancellation
