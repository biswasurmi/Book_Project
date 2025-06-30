package middleware

import (
	"net/http"
	"github.com/go-chi/jwtauth/v5"
)

func JWTAuth(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, "Invalid or missing JWT token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}