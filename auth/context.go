package auth

import (
	"context"
	"net/http"
)

type AuthInfo struct {
	UserID uint
	Token  string
}

type contextKey string

const authContextKey = contextKey("auth")

// Store AuthInfo in context
func WithAuthInfo(ctx context.Context, info AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey, info)
}

// Retrieve AuthInfo from context
func GetAuthInfo(r *http.Request) (AuthInfo, bool) {
	info, ok := r.Context().Value(authContextKey).(AuthInfo)
	return info, ok
}

func AuthHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
}
