package securityheaders

import "net/http"

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		wtr.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		wtr.Header().Set("X-Content-Type-Options", "nosniff")
		wtr.Header().Set("X-Frame-Options", "DENY")
		wtr.Header().Set("Referrer-Policy", "no-referrer")
		wtr.Header().Set("Permissions-Policy", "geolocation=(), camera=()")
		wtr.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(wtr, req)
	})
}