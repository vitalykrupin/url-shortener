package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type GetAllUserURLs struct {
	BaseHandler
}

type getUserURLsResponseUnit struct {
	Alias       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewGetAllUserURLs(app *app.App) *GetAllUserURLs {
	return &GetAllUserURLs{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *GetAllUserURLs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()

	if req.Method != http.MethodGet {
		log.Println("Only GET requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userUUIDAny := ctx.Value(middleware.UserIDKey)
	if userUUIDAny == nil || userUUIDAny == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID := userUUIDAny.(string)
	urls, err := storage.Store.GetUserURLs(ctx, userUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		log.Println(userUUIDAny)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := make([]getUserURLsResponseUnit, 0, len(urls))
	for alias, originalURL := range urls {
		result = append(result, getUserURLsResponseUnit{
			Alias:       handler.app.Config.ResponseAddress + "/" + string(alias),
			OriginalURL: string(originalURL),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
