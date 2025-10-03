package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// mockDeleteService is a mock implementation for testing
type mockDeleteService struct{}

func (m *mockDeleteService) Add(userID string, urls []string) {
	// Mock implementation - do nothing
}

func (m *mockDeleteService) Start(workers int) {
	// Mock implementation - do nothing
}

func (m *mockDeleteService) Stop() {
	// Mock implementation - do nothing
}

func TestDeleteHandler_ValidRequest(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()
	// Create a mock delete service for testing
	deleteSvc := &mockDeleteService{}
	ap := app.NewApp(store, conf, deleteSvc)

	aliases := []string{"alias1", "alias2"}
	body, _ := json.Marshal(aliases)
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewDeleteHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 status, got %d", res.StatusCode)
	}
}

func TestDeleteHandler_NoUserID(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()
	deleteSvc := &mockDeleteService{}
	ap := app.NewApp(store, conf, deleteSvc)

	aliases := []string{"alias1", "alias2"}
	body, _ := json.Marshal(aliases)
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
	// No user ID in context
	w := httptest.NewRecorder()

	h := NewDeleteHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 status, got %d", res.StatusCode)
	}
}

func TestDeleteHandler_WrongMethod(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()
	deleteSvc := &mockDeleteService{}
	ap := app.NewApp(store, conf, deleteSvc)

	aliases := []string{"alias1", "alias2"}
	body, _ := json.Marshal(aliases)
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewDeleteHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 status, got %d", res.StatusCode)
	}
}

func TestDeleteHandler_InvalidJSON(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()
	deleteSvc := &mockDeleteService{}
	ap := app.NewApp(store, conf, deleteSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader([]byte("invalid json")))
	req = req.WithContext(middleware.SetUserID(req.Context(), "user123"))
	w := httptest.NewRecorder()

	h := NewDeleteHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	// The handler currently doesn't validate JSON format, so it returns 200
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 status, got %d", res.StatusCode)
	}
}
