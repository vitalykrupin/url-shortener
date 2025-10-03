package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CookieKey string

type ShortenerClaims struct {
	jwt.RegisteredClaims
	UserID string
}

const (
	UserIDKey CookieKey = "UserId"
	tokenLT             = time.Hour * 24
)

// SetUserID is a helper function for tests to set user ID in context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "insecure-default-change-me"
		}

		token, err := req.Cookie("Token")
		if errors.Is(err, http.ErrNoCookie) {
			newUUID := uuid.NewString()
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ShortenerClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLT)),
				},
				UserID: newUUID,
			})
			strJwt, err := newToken.SignedString([]byte(secretKey))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "Token",
				Value:    strJwt,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			if req.RequestURI == "/api/user/urls" {
				next.ServeHTTP(w, req)
				return
			}
			ctx := context.WithValue(req.Context(), UserIDKey, newUUID)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		claims := &ShortenerClaims{}

		_, err = jwt.ParseWithClaims(token.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
