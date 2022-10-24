// Package middleware
package middleware

import (
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
)

// MiddlewareFunc is contract for middleware and must implement this type for http if need middleware http request
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, conf *appctx.Config) bool

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(w http.ResponseWriter, r *http.Request, conf *appctx.Config, mfs []MiddlewareFunc) bool {
	for _, mf := range mfs {
		return mf(w, r, conf)
	}

	return true
}
