package middleware

import (
	"komodo-forge-sdk-go/http/middleware/auth"
	clienttype "komodo-forge-sdk-go/http/middleware/client-type"

	// "komodo-forge-sdk-go/http/middleware/context"
	"komodo-forge-sdk-go/http/middleware/chain"
	"komodo-forge-sdk-go/http/middleware/cors"
	"komodo-forge-sdk-go/http/middleware/csrf"
	"komodo-forge-sdk-go/http/middleware/idempotency"
	ipaccess "komodo-forge-sdk-go/http/middleware/ip-access"
	"komodo-forge-sdk-go/http/middleware/normalization"
	ratelimiter "komodo-forge-sdk-go/http/middleware/rate-limiter"
	"komodo-forge-sdk-go/http/middleware/redaction"
	requestid "komodo-forge-sdk-go/http/middleware/request-id"
	rulevalidation "komodo-forge-sdk-go/http/middleware/rule-validation"
	"komodo-forge-sdk-go/http/middleware/sanitization"
	"komodo-forge-sdk-go/http/middleware/scope"
	securityheaders "komodo-forge-sdk-go/http/middleware/security-headers"
	telemetry "komodo-forge-sdk-go/http/middleware/telemetry"
)

var (
	AuthMiddleware = auth.AuthMiddleware
	Chain = chain.Chain
	ClientTypeMiddleware = clienttype.ClientTypeMiddleware
	// ContextMiddleware = context.ContextMiddleware
	CORSMiddleware = cors.CORSMiddleware
	CSRFMiddleware = csrf.CSRFMiddleware
	IdempotencyMiddleware = idempotency.IdempotencyMiddleware
	IPAccessMiddleware = ipaccess.IPAccessMiddleware
	NormalizationMiddleware = normalization.NormalizationMiddleware
	RateLimiterMiddleware = ratelimiter.RateLimiterMiddleware
	RedactionMiddleware = redaction.RedactionMiddleware
	RequestIDMiddleware = requestid.RequestIDMiddleware
	RuleValidationMiddleware = rulevalidation.RuleValidationMiddleware
	SanitizationMiddleware = sanitization.SanitizationMiddleware
	SecurityHeadersMiddleware = securityheaders.SecurityHeadersMiddleware
	ScopeMiddleware = scope.RequireServiceScope
	TelemetryMiddleware = telemetry.TelemetryMiddleware
)
