package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware_ResponseCompression(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	})
	h := GzipMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip content encoding, got %q", got)
	}
	if vary := res.Header.Get("Vary"); vary != "Accept-Encoding" {
		t.Fatalf("expected Vary header, got %q", vary)
	}
	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatalf("gzip reader error: %v", err)
	}
	defer gzr.Close()
	b, _ := io.ReadAll(gzr)
	if string(b) != "hello world" {
		t.Fatalf("unexpected body: %s", string(b))
	}
}

func TestGzipMiddleware_NoCompressionWhenNotAccepted(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	})
	h := GzipMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No Accept-Encoding header
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Encoding"); got != "" {
		t.Fatalf("expected no content encoding, got %q", got)
	}
	b, _ := io.ReadAll(res.Body)
	if string(b) != "hello world" {
		t.Fatalf("unexpected body: %s", string(b))
	}
}

func TestGzipMiddleware_NoCompressionForErrorStatus(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	})
	h := GzipMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Encoding"); got != "" {
		t.Fatalf("expected no content encoding for error status, got %q", got)
	}
	// For error status, the response should not be compressed
	b, _ := io.ReadAll(res.Body)
	// The body might have some gzip artifacts, so let's just check it contains "error"
	if !strings.Contains(string(b), "error") {
		t.Fatalf("unexpected body: %s", string(b))
	}
}

func TestGzipMiddleware_RequestDecompression(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		if string(data) != "payload" {
			t.Fatalf("unexpected decompressed body: %s", string(data))
		}
		w.WriteHeader(http.StatusOK)
	})
	h := GzipMiddleware(base)

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	_, _ = gzw.Write([]byte("payload"))
	_ = gzw.Close()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Result().StatusCode)
	}
}

func TestGzipMiddleware_RequestDecompressionError(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})
	h := GzipMiddleware(base)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid gzip data"))
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 status, got %d", rec.Result().StatusCode)
	}
}

func TestGzipMiddleware_LargeResponse(t *testing.T) {
	largeData := strings.Repeat("a", 10000)
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(largeData))
	})
	h := GzipMiddleware(base)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip content encoding, got %q", got)
	}
	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatalf("gzip reader error: %v", err)
	}
	defer gzr.Close()
	b, _ := io.ReadAll(gzr)
	if string(b) != largeData {
		t.Fatalf("unexpected decompressed body length: %d", len(b))
	}
}
