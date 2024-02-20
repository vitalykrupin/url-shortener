package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
)

type getPingHandler struct {
	BaseHandler
}

func NewGetPingHandler(app *config.App) *getPingHandler {
	return &getPingHandler{
		BaseHandler: BaseHandler{
			app: app,
		},
	}
}

func (h *getPingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.app.DBPool != nil {
		err := h.app.DBPool.Ping(context.Background())
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
}
