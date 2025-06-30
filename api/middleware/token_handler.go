package middleware

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/go-chi/jwtauth/v5"
)

func GetTokenHandler(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth, authEnabled bool) {
	if authEnabled {
		username, password, ok := r.BasicAuth()
		if !ok || username != "admin" || password != "admin123" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	exp := time.Now().Add(100000 * time.Minute).Unix()
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id":  123,
		"username": "admin",
		"exp":      exp,
	})
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}