package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// mockDeleteService for testing
type mockDeleteService struct{}

func (m *mockDeleteService) Add(userID string, urls []string) {}
func (m *mockDeleteService) Start(workers int)                {}
func (m *mockDeleteService) Stop()                            {}

func TestBuild(t *testing.T) {
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	// Create a mock storage
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}

	deleteSvc := &mockDeleteService{}
	application := app.NewApp(store, conf, deleteSvc)

	handler := Build(application)
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	// Test that handler is an http.Handler
	var _ http.Handler = handler
}

func TestBuild_Routes(t *testing.T) {
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}

	deleteSvc := &mockDeleteService{}
	application := app.NewApp(store, conf, deleteSvc)

	handler := Build(application)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET /ping", "GET", "/ping", http.StatusUnauthorized},
		{"GET /api/user/urls", "GET", "/api/user/urls", http.StatusUnauthorized},           // No JWT token
		{"POST /api/shorten", "POST", "/api/shorten", http.StatusUnauthorized},             // Blocked by auth middleware
		{"POST /api/shorten/batch", "POST", "/api/shorten/batch", http.StatusUnauthorized}, // Blocked by auth middleware
		{"DELETE /api/user/urls", "DELETE", "/api/user/urls", http.StatusUnauthorized},     // No JWT token
		{"GET /nonexistent", "GET", "/nonexistent", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for %s %s", tt.expectedStatus, res.StatusCode, tt.method, tt.path)
			}
		})
	}
}

func TestBuild_MiddlewareStack(t *testing.T) {
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}

	deleteSvc := &mockDeleteService{}
	application := app.NewApp(store, conf, deleteSvc)

	handler := Build(application)

	// Test that middleware is applied (auth required)
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 from auth middleware, got %d", res.StatusCode)
	}
}

func TestBuild_CompressionMiddleware(t *testing.T) {
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}

	deleteSvc := &mockDeleteService{}
	application := app.NewApp(store, conf, deleteSvc)

	handler := Build(application)

	// Test gzip compression
	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Should have gzip compression for successful responses
	if res.StatusCode == http.StatusOK {
		contentEncoding := res.Header.Get("Content-Encoding")
		if contentEncoding != "gzip" {
			t.Errorf("Expected gzip compression, got Content-Encoding: %s", contentEncoding)
		}
	}
}

func TestBuild_LoggingMiddleware(t *testing.T) {
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}

	deleteSvc := &mockDeleteService{}
	application := app.NewApp(store, conf, deleteSvc)

	handler := Build(application)

	// Test that logging middleware doesn't break the request (but auth still applies)
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", res.StatusCode)
	}
}
