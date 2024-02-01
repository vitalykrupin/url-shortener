package utils

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RandomString(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]rune, size)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = chars[rnd.Intn(len(chars))]
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
