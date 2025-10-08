// Package handlers provides HTTP request handlers
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
	"github.com/vitalykrupin/url-shortener/internal/app/utils"
)

// postBatchRequestUnit represents a single unit in the batch request
type postBatchRequestUnit struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

// postBatchResponseUnit represents a single unit in the batch response
type postBatchResponseUnit struct {
	CorrelationID string `json:"correlation_id"`
	Alias         string `json:"short_url"`
}

// PostBatchHandler handles POST requests for batch URL creation
type PostBatchHandler struct {
	BaseHandler
}

// NewPostBatchHandler is the constructor for PostBatchHandler
func NewPostBatchHandler(app *app.App) *PostBatchHandler {
	return &PostBatchHandler{
		BaseHandler: BaseHandler{app},
	}
}

// ServeHTTP handles the HTTP request for batch URL creation
func (handler *PostBatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		jsonReq []postBatchRequestUnit
		resp    []postBatchResponseUnit
	)
	err := json.NewDecoder(req.Body).Decode(&jsonReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	batch := make(map[storage.Alias]storage.OriginalURL)
	for _, v := range jsonReq {
		if alias, err := handler.app.Store.GetAlias(ctx, storage.OriginalURL(v.URL)); err == nil {
			resp = append(resp, postBatchResponseUnit{CorrelationID: v.CorrelationID, Alias: handler.app.Config.ResponseAddress + "/" + string(alias)})
		} else {
			alias := utils.RandomString(aliasSize)
			batch[storage.Alias(alias)] = storage.OriginalURL(v.URL)
			resp = append(resp, postBatchResponseUnit{CorrelationID: v.CorrelationID, Alias: handler.app.Config.ResponseAddress + "/" + alias})
		}
	}
	if len(batch) > 0 {
		if err := handler.app.Store.Add(ctx, batch); err != nil {
			log.Println("Can not add note to database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
