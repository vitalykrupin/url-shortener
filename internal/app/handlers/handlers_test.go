package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func AddChiContext(r *http.Request, params map[string]string) *http.Request {
	c := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, c))
	for key, val := range params {
		c.URLParams.Add(key, val)
	}

	return req
}

func TestPostHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code      int
		response    string
		contentType string
	}

	store := storage.NewStorage()
	store.FullURLKeysMap["https://yandex.ru"] = "abcABC"
	store.AliasKeysMap["abcABC"] = "https://yandex.ru"
	conf := &config.Config{}
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

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
				response: "abcABC",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPostHandler(store, conf)
			handler.ServeHTTP(tt.args.w, tt.args.req)
			res := tt.args.w.Result()
			defer res.Body.Close()
			require.Equal(t, tt.want.code, res.StatusCode)
			tt.args.w.Flush()
			if res.StatusCode != http.StatusBadRequest {
				bodyStr, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				require.NotEmpty(t, bodyStr)
			}
		})
	}
}

func TestGetHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code     int
		location string
	}

	store := storage.NewStorage()
	store.FullURLKeysMap["https://yandex.ru"] = "abcABC"
	store.AliasKeysMap["abcABC"] = "https://yandex.ru"
	conf := &config.Config{}
	conf.ServerAddress = "localhost:8080"
	conf.ResponseAddress = "http://localhost:8080"

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive GET handler test",
			args: args{
				w: httptest.NewRecorder(),
				req: AddChiContext(httptest.NewRequest(http.MethodGet, "/abcABC", nil), map[string]string{"id": "abcABC"}),
			},
			want: want{
				code:   http.StatusTemporaryRedirect,
				location: "https://yandex.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGetHandler(store, conf)
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
