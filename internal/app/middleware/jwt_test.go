package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtAuthorization_NoToken_Unauthorized(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := JWTMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/some", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}
}

func TestJwtAuthorization_InvalidToken_Unauthorized(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called with invalid token")
	})
	h := JWTMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/some", nil)
	req.AddCookie(&http.Cookie{Name: "Token", Value: "invalid.token.here"})
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 status, got %d", res.StatusCode)
	}
}

// keep one invalid test above; remove duplicates

func TestJwtAuthorization_ProtectedEndpoint_Unauthorized(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := JWTMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 status for /api/user/urls, got %d", res.StatusCode)
	}
}
