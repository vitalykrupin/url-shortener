package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtAuthorization_SetsCookieAndContext(t *testing.T) {
	var gotUser string
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(UserIDKey)
		if val == nil {
			t.Fatalf("expected user id in context")
		}
		gotUser, _ = val.(string)
		w.WriteHeader(http.StatusOK)
	})
	h := JwtAuthorization(base)

	req := httptest.NewRequest(http.MethodGet, "/some", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	if gotUser == "" {
		t.Fatalf("expected non-empty user id")
	}
	found := false
	for _, c := range res.Cookies() {
		if c.Name == "Token" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected jwt cookie set")
	}
}

func TestJwtAuthorization_ValidExistingToken(t *testing.T) {
	var gotUser string
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(UserIDKey)
		if val == nil {
			t.Fatalf("expected user id in context")
		}
		gotUser, _ = val.(string)
		w.WriteHeader(http.StatusOK)
	})
	h := JwtAuthorization(base)

	// First request to get a token
	req1 := httptest.NewRequest(http.MethodGet, "/some", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)
	res1 := rec1.Result()
	defer res1.Body.Close()

	// Extract token from first response
	var token string
	for _, c := range res1.Cookies() {
		if c.Name == "Token" {
			token = c.Value
			break
		}
	}
	if token == "" {
		t.Fatal("expected token in first response")
	}

	// Second request with existing token
	req2 := httptest.NewRequest(http.MethodGet, "/some", nil)
	req2.AddCookie(&http.Cookie{Name: "Token", Value: token})
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	res2 := rec2.Result()
	defer res2.Body.Close()

	if res2.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", res2.StatusCode)
	}
	if gotUser == "" {
		t.Fatalf("expected non-empty user id")
	}
}

func TestJwtAuthorization_InvalidToken(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called with invalid token")
	})
	h := JwtAuthorization(base)

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

func TestJwtAuthorization_SpecialRouteBehavior(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For /api/user/urls, middleware should call handler even without user context
		w.WriteHeader(http.StatusOK)
	})
	h := JwtAuthorization(base)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 status for /api/user/urls, got %d", res.StatusCode)
	}
}
