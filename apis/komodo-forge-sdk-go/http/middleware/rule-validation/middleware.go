package rulevalidation

import (
	"fmt"
	httpErr "komodo-forge-sdk-go/http/errors"
	evalRules "komodo-forge-sdk-go/http/rules"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
)

// Enforces request validation rules based on predefined configurations.
func RuleValidationMiddleware(next http.Handler) http.Handler {
	// Ensure config is loaded
	if !evalRules.LoadConfig() {
		logger.Error("validation rules failed to load", fmt.Errorf("failed to load validation rules"))
	}

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if rule := evalRules.GetRule(req.URL.Path, req.Method); rule != nil {
			if !evalRules.IsRuleValid(req, rule) {
				logger.Error("request does not comply with validation rule", fmt.Errorf("validation rule failed: %v", rule))
				httpErr.SendError(
					wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("request contents invalid"),
				)
				return
			}
		} else {
			logger.Error("no validation rule found", fmt.Errorf("no validation rule found for path: %s and method: %s", req.URL.Path, req.Method))
			httpErr.SendError(
				wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to validate request"),
			)
			return
		}
		next.ServeHTTP(wtr, req)
	})
}
