// Package middleware provides HTTP middleware functions for the main application
package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ContextKey represents the context key type
type ContextKey string

const (
	// UserIDKey is the key for user ID in context
	UserIDKey ContextKey = "user_id"
	
	// DefaultAuthServerURL is the default URL for auth service
	DefaultAuthServerURL = "http://auth:8082"
)

// SetUserID is a helper function for tests to set user ID in context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// ExternalJwtAuthorization provides JWT authorization middleware for main application
// It validates tokens by calling the external auth service
func ExternalJwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get auth service URL from environment or use default
		authServerURL := os.Getenv("AUTH_SERVER_URL")
		if authServerURL == "" {
			authServerURL = DefaultAuthServerURL
		}
		
		fmt.Printf("Auth server URL: %s\n", authServerURL)
		
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("Authorization header: %s\n", authHeader)
		
		if authHeader == "" {
			// Try to get token from cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			authHeader = "Bearer " + cookie.Value
		}
		
		// Parse token
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		
		fmt.Printf("Token string: %s\n", tokenString)
		
		// Create request to auth service to validate token
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		
		// Prepare request to auth service
		req, err := http.NewRequest("GET", authServerURL+"/api/auth/profile", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		// Add authorization header
		req.Header.Set("Authorization", "Bearer "+tokenString)
		
		// Send request to auth service
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error calling auth service: %v\n", err)
			http.Error(w, "Auth service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		
		fmt.Printf("Auth service response status: %d\n", resp.StatusCode)
		
		// Check response status
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusUnauthorized {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Auth service error", http.StatusServiceUnavailable)
			return
		}
		
		// Extract user ID from response
		// Since the auth service returns "User ID: {id}", we need to parse it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}
		
		responseBody := string(body)
		fmt.Printf("Auth service response body: %s\n", responseBody)
		
		userID := strings.TrimSpace(strings.TrimPrefix(responseBody, "User ID: "))
		fmt.Printf("Extracted user ID: %s\n", userID)
		
		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
