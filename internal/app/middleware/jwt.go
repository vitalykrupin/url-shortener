// Package middleware provides HTTP middleware functions for the main application
package middleware

import (
	"context"
	"net/http"

	auth "github.com/vitalykrupin/auth-service/pkg/auth"
)

// ContextKey is an alias to the auth-service context key type.
type ContextKey = auth.ContextKey

// UserIDKey re-exports the user id context key from auth-service.
const UserIDKey = auth.UserIDKey

// SetUserID is a helper function for tests to set user ID in context.
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// JWTMiddleware delegates JWT validation to the auth-service package middleware.
func JWTMiddleware(next http.Handler) http.Handler { return auth.JWTMiddleware(next) }
