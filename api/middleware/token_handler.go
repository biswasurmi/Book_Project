package middleware

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/biswasurmi/book-cli/service"
	"github.com/go-chi/jwtauth/v5"
)

func GetTokenHandler(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth, authEnabled bool, userService service.UserService) {
    if authEnabled {
        email, password, ok := r.BasicAuth() // Changed 'username' to 'email'
        if !ok {
            w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        user, err := userService.Authenticate(email, password)
        if err != nil {
            w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        _, tokenString, err := tokenAuth.Encode(map[string]interface{}{
            "user_id": user.ID,
            "email":   user.Email,
            "exp":     time.Now().Add(time.Hour * 24).Unix(),
        })
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "token": tokenString,
        })
        return
    }

    // If auth is disabled, return a default token
    _, tokenString, err := tokenAuth.Encode(map[string]interface{}{
        "user_id": 0,
        "email":   "test@example.com",
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
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