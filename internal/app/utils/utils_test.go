package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 0", 0},
		{"size 1", 1},
		{"size 7", 7},
		{"size 10", 10},
		{"size 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.size)
			if len(result) != tt.size {
				t.Errorf("RandomString(%d) = length %d, want %d", tt.size, len(result), tt.size)
			}

			// Check that all characters are from the expected set
			expectedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
			for _, char := range result {
				found := false
				for _, expected := range expectedChars {
					if char == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("RandomString() contains unexpected character: %c", char)
				}
			}
		})
	}
}

func TestRandomString_Uniqueness(t *testing.T) {
	// Test that multiple calls produce different results (very high probability)
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result := RandomString(10)
		if results[result] {
			t.Errorf("RandomString() produced duplicate result: %s", result)
		}
		results[result] = true
	}
}

func TestAddChiContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	params := map[string]string{
		"id":   "123",
		"name": "test",
	}

	result := AddChiContext(req, params)

	// Check that context was added
	ctx := result.Context()
	routeCtx := chi.RouteContext(ctx)
	if routeCtx == nil {
		t.Fatal("Expected route context to be set")
	}

	// Check that parameters were added
	if routeCtx.URLParam("id") != "123" {
		t.Errorf("Expected id parameter to be '123', got '%s'", routeCtx.URLParam("id"))
	}
	if routeCtx.URLParam("name") != "test" {
		t.Errorf("Expected name parameter to be 'test', got '%s'", routeCtx.URLParam("name"))
	}
}

func TestAddChiContext_EmptyParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	params := map[string]string{}

	result := AddChiContext(req, params)

	// Check that context was still added
	ctx := result.Context()
	routeCtx := chi.RouteContext(ctx)
	if routeCtx == nil {
		t.Fatal("Expected route context to be set even with empty params")
	}
}

func TestWithHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	key := "Content-Type"
	value := "application/json"

	result := WithHeader(req, key, value)

	// Check that header was added
	if result.Header.Get(key) != value {
		t.Errorf("Expected header %s to be '%s', got '%s'", key, value, result.Header.Get(key))
	}
}

func TestWithHeader_MultipleValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	key := "Accept"
	value1 := "application/json"
	value2 := "text/plain"

	result1 := WithHeader(req, key, value1)
	result2 := WithHeader(result1, key, value2)

	// Check that both values were added
	headers := result2.Header[key]
	if len(headers) != 2 {
		t.Errorf("Expected 2 header values, got %d", len(headers))
	}
	if headers[0] != value1 {
		t.Errorf("Expected first header value to be '%s', got '%s'", value1, headers[0])
	}
	if headers[1] != value2 {
		t.Errorf("Expected second header value to be '%s', got '%s'", value2, headers[1])
	}
}
