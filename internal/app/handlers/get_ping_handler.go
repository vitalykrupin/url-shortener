package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type GetPingHandler struct {
	BaseHandler
}

func NewGetPingHandler(app *app.App) *GetPingHandler {
	return &GetPingHandler{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *GetPingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()

	if storage.Store != nil {
		err := storage.Store.PingStorage(ctx)
		if err != nil {
			log.Println("Can not connect to database")
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
