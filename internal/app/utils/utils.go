package utils

import (
	"context"
	"crypto/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RandomString(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	if size <= 0 {
		return ""
	}
	result := make([]rune, size)
	// Read crypto-strong random bytes and map into chars range
	randomBytes := make([]byte, size)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback: deterministic but safe empty string on error
		return ""
	}
	for i := 0; i < size; i++ {
		result[i] = chars[int(randomBytes[i])%len(chars)]
	}
	return string(result)
}

func AddChiContext(r *http.Request, params map[string]string) *http.Request {
	c := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, c))
	for key, val := range params {
		c.URLParams.Add(key, val)
	}

	return req
}

func WithHeader(req *http.Request, key string, value string) *http.Request {
	req.Header.Add(key, value)
	return req
}
