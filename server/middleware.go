package server

import (
	"context"
	"github.com/golang-jwt/jwt"
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

func (t *TodoServer) Authorize() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := getToken(r)
			token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(t.config.SecretKey), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(*model.Claims)
			if !ok || !token.Valid {
				http.Error(w, "invalid token: authentication failed", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userid", claims.UserId)
			r = r.WithContext(ctx)

			f(w, r)
		}
	}
}

func (t *TodoServer) Log() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func (t *TodoServer) SetContentType() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, r)
		}
	}
}
