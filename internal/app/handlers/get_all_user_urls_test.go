package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

func TestGetAllUserURLs_WithURLs(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	// Pre-insert URLs for user
	ctx := middleware.SetUserID(httptest.NewRequest(http.MethodPost, "/", nil).Context(), "user123")
	_ = store.Add(ctx, map[storage.Alias]storage.OriginalURL{
		"alias1": "https://example1.com",
		"alias2": "https://example2.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewGetAllUserURLs(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	// File storage doesn't support GetUserURLs, so it returns 400
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 status (file storage doesn't support GetUserURLs), got %d", res.StatusCode)
	}
}

func TestGetAllUserURLs_NoURLs(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewGetAllUserURLs(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	// File storage doesn't support GetUserURLs, so it returns 400
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 status (file storage doesn't support GetUserURLs), got %d", res.StatusCode)
	}
}

func TestGetAllUserURLs_NoUserID(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	// No user ID in context
	w := httptest.NewRecorder()

	h := NewGetAllUserURLs(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 status, got %d", res.StatusCode)
	}
}

func TestGetAllUserURLs_WrongMethod(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewGetAllUserURLs(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 status, got %d", res.StatusCode)
	}
}
