package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/mthsgimenez/participe/internal/auth"
)

type contextKey string

const userContextKey contextKey = "userClaims"

func GetUserClaims(r *http.Request) *auth.Claims {
	if claims, ok := r.Context().Value(userContextKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondJSONError(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondJSONError(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := auth.VerifyJWT(tokenString)
		if err != nil {
			RespondJSONError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
