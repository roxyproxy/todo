package httpsrv

import (
	"context"
	"fmt"
	"net/http"
	"todo/model"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func (t *Server) Authorize() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := getToken(r)
			claims, err := t.service.ValidateToken(tokenString)

			if err != nil {
				t.handleError(fmt.Errorf("%q: %q: %w", "authentication failed.", err, model.ErrUnauthorized), w)
				return
			}

			ctx := context.WithValue(r.Context(), model.KeyUserId("userid"), claims.UserId)
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

func (t *Server) Log() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func (t *Server) SetContentType() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, r)
		}
	}
}
