package headerEval

import (
	"komodo-forge-sdk-go/config"
	"komodo-forge-sdk-go/crypto/jwt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ValidateHeaderValue runs lightweight validation for known header names.
func ValidateHeaderValue(hdr string, req *http.Request) (bool, error) {
	val := req.Header.Get(hdr)
	switch strings.ToLower(hdr) {
		case "access-control-allow-origin":
			return isValidCORS(val), nil
		case "authorization":
			return jwt.ValidateToken(val)
		case "cache-control":
			return isValidCacheControl(val), nil
		case "cookie":
			return isValidCookie(val), nil
		case "content-type", "accept":
			return isValidContentAcceptType(val), nil
		case "content-length":
			return isValidContentLength(val), nil
		case "idempotency-key":
			return isValidIdempotencyKey(val), nil
		case "referer", "referrer":
			return isValidReferer(val), nil
		case "user-agent":
			return isValidUserAgent(val), nil
		case "x-csrf-token":
			return isValidCSRF(val), nil
		case "x-requested-by":
			return isValidRequestedBy(val), nil
		default:
			return val != "", nil
	}
}

// ================ Helper validation functions ================

func isValidContentAcceptType(s string) bool {
	return strings.HasPrefix(s, "application/json") ||
		strings.HasPrefix(s, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(s, "multipart/form-data")
}

func isValidContentLength(s string) bool {
	if s == "" { return false }

	val, err := strconv.Atoi(s)
	if err != nil { return false }

	max := (func() int {
		val := config.GetConfigValue("MAX_CONTENT_LENGTH")
		num, err := strconv.Atoi(val)
		if val == "" || err != nil { return 4096 }
		return num
	})()

	return val > 0 && val <= max
}

func isValidCookie(s string) bool {
	// TODO: Implement cookie validation logic (e.g., parse, check signature)
	return s != ""
}

func isValidUserAgent(s string) bool {
	if s == "" { return false }
	s = strings.TrimSpace(s)
	if len(s) > 256 { return false } // max length
	re := regexp.MustCompile(`^[A-Za-z0-9\-\._ /(),:;]+$`)
	return re.MatchString(s)
}

func isValidReferer(s string) bool {
	re := regexp.MustCompile(`^https?://[A-Za-z0-9\-.%]+(?::\d{1,5})?(?:/.*)?$`)
	return re.MatchString(strings.TrimSpace(s))
}

func isValidCacheControl(s string) bool {
	return s == "no-cache" || s == "no-store" || s == "must-revalidate"
}

func isValidRequestedBy(s string) bool {
  return s != "" && len(s) <= 64 && regexp.MustCompile(`^[A-Za-z0-9_\-/]+$`).MatchString(s)
}

func isValidIdempotencyKey(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-]{8,64}$`).MatchString(s)
}

func isValidCSRF(s string) bool {
	if s == "" { return false }
	// TODO - implement CSRF token validation logic
	return true
}

func isValidCORS(s string) bool {
	if s == "" { return false }
	if s == "*" { return true }
	re := regexp.MustCompile(`^https?://[A-Za-z0-9\-.%]+(?::\d{1,5})?(?:/.*)?$`)
	return re.MatchString(strings.TrimSpace(s))
}
