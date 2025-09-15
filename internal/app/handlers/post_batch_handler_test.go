package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

type batchReq struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

type batchResp struct {
	CorrelationID string `json:"correlation_id"`
	Alias         string `json:"short_url"`
}

func TestPostBatchHandler_MixedExistingAndNew(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	// Pre-insert existing URL
	_ = store.Add(httptest.NewRequest(http.MethodPost, "/", nil).Context(), map[storage.Alias]storage.OriginalURL{"exist01": "https://exist"})

	body, _ := json.Marshal([]batchReq{{CorrelationID: "1", URL: "https://exist"}, {CorrelationID: "2", URL: "https://new"}})
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	var resp []batchResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
	if resp[0].CorrelationID != "1" || resp[0].Alias == "" {
		t.Fatalf("unexpected first resp: %+v", resp[0])
	}
	if resp[1].CorrelationID != "2" || resp[1].Alias == "" {
		t.Fatalf("unexpected second resp: %+v", resp[1])
	}
}

func TestPostBatchHandler_AllExisting(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	// Pre-insert existing URLs
	_ = store.Add(httptest.NewRequest(http.MethodPost, "/", nil).Context(), map[storage.Alias]storage.OriginalURL{
		"exist01": "https://exist1",
		"exist02": "https://exist2",
	})

	body, _ := json.Marshal([]batchReq{
		{CorrelationID: "1", URL: "https://exist1"},
		{CorrelationID: "2", URL: "https://exist2"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	var resp []batchResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
}

func TestPostBatchHandler_AllNew(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	body, _ := json.Marshal([]batchReq{
		{CorrelationID: "1", URL: "https://new1"},
		{CorrelationID: "2", URL: "https://new2"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	var resp []batchResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
}

func TestPostBatchHandler_InvalidJSON(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 status, got %d", res.StatusCode)
	}
}

func TestPostBatchHandler_WrongContentType(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	body, _ := json.Marshal([]batchReq{{CorrelationID: "1", URL: "https://test"}})
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 status, got %d", res.StatusCode)
	}
}

func TestPostBatchHandler_WrongMethod(t *testing.T) {
	conf := config.NewConfig()
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = filepath.Join(t.TempDir(), "testfile.json")
	store, err := storage.NewStorage(conf)
	if err != nil {
		t.Fatal(err)
	}
	ap := app.NewApp(store, conf, nil)

	body, _ := json.Marshal([]batchReq{{CorrelationID: "1", URL: "https://test"}})
	req := httptest.NewRequest(http.MethodGet, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := NewPostBatchHandler(ap)
	h.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 status, got %d", res.StatusCode)
	}
}
