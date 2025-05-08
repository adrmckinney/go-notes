package middleware

import (
	"net/http"
	"strings"

	"github.com/adrmckinney/go-notes/auth"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userId"

func AuthMiddleware(jwtKey []byte, tokenRepo repos.UserTokenRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !tokenRepo.TokenExists(tokenString) {
				http.Error(w, "Token revoked", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "Missing user ID", http.StatusUnauthorized)
				return
			}
			userID := uint(userIDFloat)

			// Store in context
			ctx := auth.WithAuthInfo(r.Context(), auth.AuthInfo{
				UserID: userID,
				Token:  tokenString,
			})
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
