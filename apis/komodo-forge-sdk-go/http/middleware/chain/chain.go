package chain

import "net/http"

// Applies middleware in order: first listed = outermost wrapper.
func Chain(f func(http.ResponseWriter, *http.Request), mw ...func(http.Handler) http.Handler) http.Handler {
	var c http.Handler = http.HandlerFunc(f)
	for i := len(mw) - 1; i >= 0; i-- { c = mw[i](c) }
	return c
}