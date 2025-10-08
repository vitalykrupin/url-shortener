package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got '%s'", contentType)
		}

		// Check Authorization header if provided
		auth := r.Header.Get("Authorization")
		if auth != "" && !strings.HasPrefix(auth, "Bearer ") {
			t.Errorf("Expected Authorization to start with 'Bearer ', got '%s'", auth)
		}

		// Read and check body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}

		expectedURL := "https://example.com"
		if string(body) != expectedURL {
			t.Errorf("Expected body '%s', got '%s'", expectedURL, string(body))
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://localhost:8080/abc123"))
	}))
	defer server.Close()

	// Test sendRequest function
	response, err := sendRequest(server.URL+"/", "https://example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}

	// Check response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	expectedResponse := "http://localhost:8080/abc123"
	if string(body) != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, string(body))
	}
}

func TestSendRequest_ServerError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Test sendRequest function
	response, err := sendRequest(server.URL+"/", "https://example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", response.StatusCode)
	}
}

func TestSendRequest_InvalidEndpoint(t *testing.T) {
	// Test with invalid endpoint
	_, err := sendRequest("http://invalid-endpoint:99999/", "https://example.com")
	if err == nil {
		t.Error("Expected error for invalid endpoint")
	}
}

func TestSendRequest_EmptyURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and check body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}

		if string(body) != "" {
			t.Errorf("Expected empty body, got '%s'", string(body))
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	// Test sendRequest function with empty URL
	response, err := sendRequest(server.URL+"/", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", response.StatusCode)
	}
}

func TestSendRequest_LongURL(t *testing.T) {
	// Create a very long URL
	longURL := "https://example.com/" + strings.Repeat("very-long-path/", 100)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and check body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}

		if string(body) != longURL {
			t.Errorf("Expected long URL, got different body")
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://localhost:8080/short"))
	}))
	defer server.Close()

	// Test sendRequest function with long URL
	response, err := sendRequest(server.URL+"/", longURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}
}

func TestSendRequest_SpecialCharacters(t *testing.T) {
	// Create URL with special characters
	specialURL := "https://example.com/path?param=value&other=тест#fragment"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and check body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}

		if string(body) != specialURL {
			t.Errorf("Expected special URL, got different body")
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://localhost:8080/special"))
	}))
	defer server.Close()

	// Test sendRequest function with special characters
	response, err := sendRequest(server.URL+"/", specialURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}
}

func TestSendRequest_Headers(t *testing.T) {
	// Create a test server that checks headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that Content-Type is set correctly
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got '%s'", contentType)
		}

		// Check that User-Agent is set (Go's default)
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("Expected User-Agent header to be set")
		}

		// Set env to simulate token and ensure Authorization is set
		if os.Getenv("AUTH_TOKEN") != "" {
			if got := r.Header.Get("Authorization"); !strings.HasPrefix(got, "Bearer ") {
				t.Errorf("Expected Authorization header with Bearer, got '%s'", got)
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Test sendRequest function
	// Provide token via env to trigger Authorization header
	os.Setenv("AUTH_TOKEN", "test-token")
	defer os.Unsetenv("AUTH_TOKEN")
	response, err := sendRequest(server.URL+"/", "https://example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}
}
