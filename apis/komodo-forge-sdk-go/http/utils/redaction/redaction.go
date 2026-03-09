package redaction

import "regexp"

var (
	keyRegex = regexp.MustCompile(`(?i)^(authorization|set-cookie|password|token|bearer|ssn|pwd|secret|api_key|cvv|card_number)$`)
	piiRegex = regexp.MustCompile(`(?i)` +
		`([a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,})|` + 	// Email
		`(\b\d{3}-\d{2}-\d{4}\b)|` +                  // SSN
		`(\b(?:\d[ -]*?){13,16}\b)|` +               	// Credit Cards
		`(Bearer\s+[a-zA-Z0-9\-\._~\+\/]+=*)`,       	// Bearer Tokens
	)
)

func RedactString(val string) string {
	if val == "" { return val }
	return piiRegex.ReplaceAllString(val, "[REDACTED]")
}

func RedactPair(key string, val any) any {
	if keyRegex.MatchString(key) { return "[REDACTED]" }
	if s, ok := val.(string); ok && s != "" {
		return piiRegex.ReplaceAllString(s, "[REDACTED]")
	}
	return val
}
