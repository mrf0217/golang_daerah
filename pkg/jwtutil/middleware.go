package jwtutil

// Request Flow Link:
// For every protected route registered in main.go, AuthMiddleware sits between the rate limiter
// and the concrete handler to ensure only requests with valid JWT tokens continue down the chain.

import (
	"net/http"
)

// AuthMiddleware checks the Authorization header using VerifyToken
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		_, err := VerifyToken(authHeader)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
