package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type GetHandler struct {
	BaseHandler
}

func NewGetHandler(store *storage.Store, config config.Config) *GetHandler {
	return &GetHandler{
		BaseHandler: BaseHandler{
			store:  store,
			config: config,
		},
	}
}

func (handler *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
	if fullURL, ok := handler.store.Store.AliasKeysMap[alias]; ok {
		w.Header().Add("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		log.Println("URL by alias " + alias + " is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
