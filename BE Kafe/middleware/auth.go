package middleware

import (
	"net/http"
	"strings"

	"be_kafe/utils"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it's Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		username, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add username to request context (optional, for future use)
		r.Header.Set("X-User", username)

		// Call next handler
		next.ServeHTTP(w, r)
	}
}
