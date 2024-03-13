package middleware

import (
	"context"
	"errors"
	"net/http"
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
	secretKey string    = "secret_key"
	UserIDKey CookieKey = "UserId"
	tokenLT             = time.Hour * 24
)

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("Token")
		if errors.Is(http.ErrNoCookie, err) {
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
				Name:  "Token",
				Value: strJwt,
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
