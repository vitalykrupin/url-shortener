package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/utils"
)

type postBatchRequestUnit struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}
type postBatchResponseUnit struct {
	CorrelationID string `json:"correlation_id"`
	Alias         string `json:"short_url"`
}

type PostBatchHandler struct {
	BaseHandler
}

func NewPostBatchHandler(app *app.App) *PostBatchHandler {
	return &PostBatchHandler{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *PostBatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		return
	}

	var (
		jsonReq []postBatchRequestUnit
		resp []postBatchResponseUnit
	)
	err := json.NewDecoder(req.Body).Decode(&jsonReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range jsonReq {
		if alias, err := handler.app.Storage.GetAlias(req.Context(), v.URL); err == nil {
			err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+alias)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			alias := utils.RandomString(aliasSize)
			if err := handler.app.Storage.Add(req.Context(), alias, v.URL); err != nil {
				log.Println("Can not add note to database")
				return
			}
			resp = append(resp, postBatchResponseUnit{CorrelationID: v.CorrelationID, Alias: handler.app.Config.ResponseAddress+"/"+alias})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
