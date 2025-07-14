package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/biswasurmi/book-cli/service"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Load environment variables once during initialization
func init() {
	if err := godotenv.Load(); err != nil {
		// Log error but continue
	}
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request, authEnabled bool, userService service.UserService) {
	if authEnabled {
		email, password, ok := r.BasicAuth()
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

		// Get the secret key from environment variable
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		if len(secretKey) == 0 {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		// Generate JWT token using golang-jwt/jwt
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		tokenString, err := token.SignedString(secretKey)
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
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = 0
	claims["email"] = "test@example.com"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}