package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type GetHandler struct {
	BaseHandler
}

func NewGetHandler(app *app.App) *GetHandler {
	return &GetHandler{
		BaseHandler: BaseHandler{app},
	}
}

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
	if URL, err := handler.app.Storage.GetURL(ctx, storage.Alias(alias)); err != nil {
		log.Println("URL by alias " + alias + " is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", string(URL))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
