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
		BaseHandler: BaseHandler{
			app: app,
		},
	}
}

func (handler *GetPingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler.app.DB != nil {
		err := handler.app.DB.PingContext(req.Context())
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("Can not connect to database")
	}
	w.WriteHeader(http.StatusInternalServerError)
}
