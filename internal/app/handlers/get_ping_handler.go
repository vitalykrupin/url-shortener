// Package handlers provides HTTP request handlers
package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener/internal/app"
)

// GetPingHandler handles GET requests for ping endpoint
type GetPingHandler struct {
	BaseHandler
}

// NewGetPingHandler is the constructor for GetPingHandler
func NewGetPingHandler(app *app.App) *GetPingHandler {
	return &GetPingHandler{
		BaseHandler: BaseHandler{app},
	}
}

// ServeHTTP handles the HTTP request for ping endpoint
func (handler *GetPingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()

	if handler.app.Store != nil {
		err := handler.app.Store.PingStorage(ctx)
		if err != nil {
			log.Println("Can not connect to database")
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
