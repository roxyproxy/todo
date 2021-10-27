package httpsrv

import (
	"context"
	"fmt"
	"net/http"

	"todo/metrics"
	"todo/model"
)

// Middleware is type of http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain function to get all handler function into a chain.
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Authorize middleware.
func (t *Server) Authorize() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := getToken(r)
			claims, err := t.service.ValidateToken(r.Context(), tokenString)
			if err != nil {
				t.handleError(fmt.Errorf("%q: %q: %w", "authentication failed.", err, model.ErrUnauthorized), w)
				return
			}

			ctx := context.WithValue(r.Context(), model.KeyUserID("userid"), claims.UserID)
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

// Log middleware.
func (t *Server) Log() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			metrics.TotalRequests.WithLabelValues(r.URL.Path).Inc()

			f(w, r)
		}
	}
}

// SetContentType middleware.
func (t *Server) SetContentType() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, r)
		}
	}
}
