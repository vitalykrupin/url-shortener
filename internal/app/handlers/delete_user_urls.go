package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener.git/internal/app/services/deleter"
)

type NewDeleteUserURLs struct {
	BaseHandler
}

func NewDeleteHandler(app *app.App) *NewDeleteUserURLs {
	return &NewDeleteUserURLs{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *NewDeleteUserURLs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		log.Println("Only DELETE requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userUUIDAny := req.Context().Value(middleware.UserIDKey)
	if userUUIDAny == nil || userUUIDAny == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID := userUUIDAny.(string)

	var urls []string
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Can not read body")
		return
	}
	if err := json.Unmarshal(body, &urls); err != nil {
		log.Println("Can not unmarshal body")
		return
	}

	ds.DelService.Add(urls, userUUID)

	w.WriteHeader(http.StatusAccepted)
}
