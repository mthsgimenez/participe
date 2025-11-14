package main

import (
	"context"
	"net/http"

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
		cookie, err := r.Cookie("jwt")
		if err != nil {
			RespondJSONError(w, "missing or invalid JWT cookie", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims, err := auth.VerifyJWT(tokenString)
		if err != nil {
			RespondJSONError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
