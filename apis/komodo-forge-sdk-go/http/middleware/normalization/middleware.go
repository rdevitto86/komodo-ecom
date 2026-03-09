package normalization

import (
	"net/http"
	"net/url"
	"strings"
)

// Normalizes request data including headers, URLs, and query parameters
func NormalizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		normalizeHeaders(req)
		normalizeURL(req)
		normalizeQueryParams(req)
		req.Method = strings.ToUpper(req.Method)
		next.ServeHTTP(wtr, req)
	})
}

// Normalizes HTTP headers
func normalizeHeaders(req *http.Request) {
	for key, values := range req.Header {
		for i, value := range values {
			req.Header[key][i] = strings.TrimSpace(value)
		}
	}

	if contentType := req.Header.Get("Content-Type"); contentType != "" {
		req.Header.Set("Content-Type", strings.ToLower(strings.TrimSpace(contentType)))
	}
	if accept := req.Header.Get("Accept"); accept != "" {
		req.Header.Set("Accept", strings.ToLower(strings.TrimSpace(accept)))
	}
	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		req.Header.Set("User-Agent", strings.TrimSpace(userAgent))
	}
}

// Normalizes URL path
func normalizeURL(req *http.Request) {
	if req.URL == nil { return }

	path := req.URL.Path
	
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	
	path = strings.ReplaceAll(path, "//", "/")
	
	req.URL.Path = path
	req.RequestURI = path
	if req.URL.RawQuery != "" {
		req.RequestURI += "?" + req.URL.RawQuery
	}
}

// Normalizes query parameters
func normalizeQueryParams(req *http.Request) {
	if req.URL == nil { return }
	normalized := url.Values{}

	for key, values := range req.URL.Query() {
		normalizedKey := strings.TrimSpace(key)
		for _, value := range values {
			normalizedValue := strings.TrimSpace(value)
			
			switch normalizedValue {
				case "true", "True", "TRUE":
					normalizedValue = "true"
				case "false", "False", "FALSE":
					normalizedValue = "false"
				case "sort", "Sort", "SORT":
					normalizedValue = "sort"
				case "asc", "Asc", "ASC":
					normalizedValue = "asc"
				case "desc", "Desc", "DESC":
					normalizedValue = "desc"
				default:
					// do nothing
			}

			normalized.Add(normalizedKey, normalizedValue)
		}
	}
	req.URL.RawQuery = normalized.Encode()
}
