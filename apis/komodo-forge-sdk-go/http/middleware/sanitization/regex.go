package sanitization

import "regexp"

var (
	sqlInjectionPattern = regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|onerror|onload|<script|</script)`)
	xssPattern = regexp.MustCompile(`(?i)(<script|</script|javascript:|onerror=|onload=|<iframe|</iframe|<object|</object|<embed|</embed)`)
	pathTraversalPattern = regexp.MustCompile(`\.\.\/|\.\.\\`)
	nullBytePattern = regexp.MustCompile(`\x00`)
)
