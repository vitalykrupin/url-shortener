package handlers

import (
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
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
	if handler.app.Storage != nil {
		err := handler.app.Storage.PingStorage(req.Context())
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("Can not connect to database")
	}
	w.WriteHeader(http.StatusInternalServerError)
}
