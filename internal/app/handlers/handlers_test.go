package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	deleteService "github.com/vitalykrupin/url-shortener.git/internal/app/services/deleter"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
	"github.com/vitalykrupin/url-shortener.git/internal/app/utils"
)

func TestGetHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code     int
		location string
	}

	conf := &config.Config{}
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = "/tmp/testfile.json"
	store := storage.NewFileStorage(conf)
	var batch = map[storage.Alias]storage.OriginalURL{
		"abcABC": "https://yandex.ru",
	}
	store.Add(context.Background(), batch)

	newApp := app.NewApp(conf, store, deleteService.NewDeleteService())

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive GET handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: utils.AddChiContext(httptest.NewRequest(http.MethodGet, "/abcABC", nil), map[string]string{idParam: "abcABC"}),
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://yandex.ru",
			},
		},
		{
			name: "negative GET handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: utils.AddChiContext(httptest.NewRequest(http.MethodGet, "/abc", nil), map[string]string{idParam: "abc"}),
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGetHandler(newApp)
			handler.ServeHTTP(tt.args.w, tt.args.req)
			res := tt.args.w.Result()
			defer res.Body.Close()
			require.Equal(t, tt.want.code, res.StatusCode)
			if res.StatusCode != http.StatusBadRequest {
				require.Equal(t, tt.want.location, res.Header.Get("Location"))
				handler.ServeHTTP(tt.args.w, tt.args.req)
				newResult := tt.args.w.Result()
				defer newResult.Body.Close()
				assert.Equal(t, res.StatusCode, newResult.StatusCode)
				assert.Equal(t, res.Header.Get("Location"), newResult.Header.Get("Location"))
			}
		})
	}
}

func TestPostHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code int
	}

	conf := &config.Config{}
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = "/tmp/testfile.json"
	store := storage.NewFileStorage(conf)
	newApp := app.NewApp(conf, store, deleteService.NewDeleteService())

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive POST handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://yandex.ru")),
			},
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPostHandler(newApp)
			handler.ServeHTTP(tt.args.w, tt.args.req)
			result := tt.args.w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want.code, result.StatusCode)
			tt.args.w.Flush()
			if result.StatusCode != http.StatusBadRequest {
				bodyStr, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				require.NotEmpty(t, bodyStr)
			}
		})
	}
}

func TestPostJSONHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code int
	}

	conf := &config.Config{}
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"
	conf.FileStorePath = "/tmp/testfile.json"
	store := storage.NewFileStorage(conf)
	newApp := app.NewApp(conf, store, deleteService.NewDeleteService())

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive POST_JSON handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: utils.WithHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("{\"url\":\"https://yandex.ru\"}")), "Content-Type", "application/json"),
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "negative POST_JSON handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: utils.WithHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("https://yandex.ru")), "Content-Type", "application/json"),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid json POST_JSON handler test",
			args: args{
				w:   httptest.NewRecorder(),
				req: utils.WithHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("{\"url\":\"https://yandex.ru\"")), "Content-Type", "application/json"),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPostHandler(newApp)
			handler.ServeHTTP(tt.args.w, tt.args.req)
			result := tt.args.w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want.code, result.StatusCode)
			tt.args.w.Flush()
			if result.StatusCode != http.StatusBadRequest {
				resp := new(postJSONResponse)
				err := json.NewDecoder(result.Body).Decode(resp)
				require.NoError(t, err)
				require.NotEmpty(t, resp.Alias)
			}
		})
	}
}
