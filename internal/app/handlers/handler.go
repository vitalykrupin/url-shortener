package handlers

import (
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type BaseHandler struct {
	http.Handler
	store  *storage.Store
	config config.Config
}
