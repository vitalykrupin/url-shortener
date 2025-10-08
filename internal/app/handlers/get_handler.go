// Package handlers provides HTTP request handlers
package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// GetHandler handles GET requests for URL redirection
type GetHandler struct {
	BaseHandler
}

// NewGetHandler is the constructor for GetHandler
func NewGetHandler(app *app.App) *GetHandler {
	return &GetHandler{
		BaseHandler: BaseHandler{app},
	}
}

// ServeHTTP handles the HTTP request for URL redirection
func (handler *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()

	if req.Method != http.MethodGet {
		log.Println("Only GET requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	alias := chi.URLParam(req, idParam)
	if alias == "" {
		log.Println("Get query require Id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if URL, err := handler.app.Store.GetURL(ctx, storage.Alias(alias)); err != nil {
		if errors.Is(err, storage.ErrDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		log.Println("URL by alias " + alias + " is not exists")
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.Header().Add("Location", string(URL))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
