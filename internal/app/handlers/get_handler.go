package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
)

type GetHandler struct {
	BaseHandler
}

func NewGetHandler(app *config.App) *GetHandler {
	return &GetHandler{
		BaseHandler: BaseHandler{
			app: app,
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
	if URL, ok := handler.app.Storage.GetURL(alias); ok {
		w.Header().Add("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		log.Println("URL by alias " + alias + " is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
