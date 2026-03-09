package scope

import (
	"net/http"
	"strings"

	ctxKeys "komodo-forge-sdk-go/http/context"
	httpErr "komodo-forge-sdk-go/http/errors"
)

// RequireServiceScope rejects requests whose JWT does not carry a service-scoped token.
// Service tokens must have at least one scope with a "svc:" prefix, issued by komodo-auth-api
// on service bootstrap — distinct from user tokens, which carry user-scoped claims.
func RequireServiceScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		scopes, ok := req.Context().Value(ctxKeys.SCOPES_KEY).([]string)
		if !ok || len(scopes) == 0 {
			httpErr.SendError(wtr, req, httpErr.Auth.InsufficientScope)
			return
		}

		for _, s := range scopes {
			if strings.HasPrefix(s, "svc:") {
				next.ServeHTTP(wtr, req)
				return
			}
		}

		httpErr.SendError(wtr, req, httpErr.Auth.InsufficientScope)
	})
}