// Package handlers provides HTTP request handlers tests
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

func TestPostHandler_ServeHTTP_TextPlain(t *testing.T) {
	// Setup
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")

	store, err := storage.NewStorage(conf)
	require.NoError(t, err)
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()

	newApp := app.NewApp(store, conf, nil)

	// Test case 1: Valid text/plain request
	t.Run("valid text/plain request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		handler := NewPostHandler(newApp)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		// Check that response body is not empty and contains expected format
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(res.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, body.String())
		assert.Contains(t, body.String(), "http://localhost:8080/")
	})

	// Test case 2: Empty body
	t.Run("empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		handler := NewPostHandler(newApp)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestPostHandler_ServeHTTP_ApplicationJSON(t *testing.T) {
	// Setup
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")

	store, err := storage.NewStorage(conf)
	require.NoError(t, err)
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()

	newApp := app.NewApp(store, conf, nil)

	// Test case 1: Valid application/json request
	t.Run("valid application/json request", func(t *testing.T) {
		jsonReq := `{"url":"https://example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := NewPostHandler(newApp)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		// Check response structure
		var resp postJSONResponse
		err := json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Alias)
		assert.Contains(t, resp.Alias, "http://localhost:8080/")
	})

	// Test case 2: Invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		jsonReq := `{"url":"https://example.com"`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := NewPostHandler(newApp)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	// Test case 3: Empty JSON body
	t.Run("empty JSON body", func(t *testing.T) {
		jsonReq := `{}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := NewPostHandler(newApp)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestPostHandler_ServeHTTP_WrongMethod(t *testing.T) {
	// Setup
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")

	store, err := storage.NewStorage(conf)
	require.NoError(t, err)
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()

	newApp := app.NewApp(store, conf, nil)

	// Test case: GET request instead of POST
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler := NewPostHandler(newApp)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestPostHandler_ServeHTTP_ExistingURL(t *testing.T) {
	// Setup
	conf := config.NewConfig()
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")

	store, err := storage.NewStorage(conf)
	require.NoError(t, err)
	defer func() {
		_ = store.CloseStorage(context.Background())
	}()

	newApp := app.NewApp(store, conf, nil)

	// First request to add URL
	req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	req1.Header.Set("Content-Type", "text/plain")
	w1 := httptest.NewRecorder()

	handler := NewPostHandler(newApp)
	handler.ServeHTTP(w1, req1)

	res1 := w1.Result()
	defer res1.Body.Close()

	assert.Equal(t, http.StatusCreated, res1.StatusCode)

	// Second request with the same URL should return conflict
	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	req2.Header.Set("Content-Type", "text/plain")
	w2 := httptest.NewRecorder()

	handler.ServeHTTP(w2, req2)

	res2 := w2.Result()
	defer res2.Body.Close()

	assert.Equal(t, http.StatusConflict, res2.StatusCode)
}

func TestParseBody_TextPlain(t *testing.T) {
	// Test case 1: Valid text/plain body
	t.Run("valid text/plain body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		req.Header.Set("Content-Type", "text/plain")

		url, err := parseBody(req)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", url)
	})

	// Test case 2: Empty text/plain body
	t.Run("empty text/plain body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		req.Header.Set("Content-Type", "text/plain")

		_, err := parseBody(req)
		assert.Error(t, err)
	})
}

func TestParseBody_ApplicationJSON(t *testing.T) {
	// Test case 1: Valid JSON body
	t.Run("valid JSON body", func(t *testing.T) {
		jsonReq := `{"url":"https://example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		url, err := parseBody(req)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", url)
	})

	// Test case 2: Invalid JSON body
	t.Run("invalid JSON body", func(t *testing.T) {
		jsonReq := `{"url":"https://example.com"`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		_, err := parseBody(req)
		assert.Error(t, err)
	})

	// Test case 3: Empty JSON body
	t.Run("empty JSON body", func(t *testing.T) {
		jsonReq := `{}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		url, err := parseBody(req)
		require.NoError(t, err)
		assert.Equal(t, "", url)
	})
}

func TestPrintResponse_TextPlain(t *testing.T) {
	// Test case 1: Text/plain response for new URL
	t.Run("text/plain response for new URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		err := printResponse(w, req, "http://localhost:8080/abc123", false)
		require.NoError(t, err)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		body := new(bytes.Buffer)
		_, err = body.ReadFrom(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/abc123", body.String())
	})

	// Test case 2: Text/plain response for existing URL
	t.Run("text/plain response for existing URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		err := printResponse(w, req, "http://localhost:8080/abc123", true)
		require.NoError(t, err)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusConflict, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		body := new(bytes.Buffer)
		_, err = body.ReadFrom(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/abc123", body.String())
	})
}

func TestPrintResponse_ApplicationJSON(t *testing.T) {
	// Test case 1: JSON response for new URL
	t.Run("JSON response for new URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		err := printResponse(w, req, "http://localhost:8080/abc123", false)
		require.NoError(t, err)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		var resp postJSONResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/abc123", resp.Alias)
	})

	// Test case 2: JSON response for existing URL
	t.Run("JSON response for existing URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		err := printResponse(w, req, "http://localhost:8080/abc123", true)
		require.NoError(t, err)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusConflict, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		var resp postJSONResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/abc123", resp.Alias)
	})
}
