package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
	"github.com/vitalykrupin/url-shortener.git/internal/app/utils"
)

type postJSONRequest struct {
	URL string `json:"url"`
}

type postJSONResponse struct {
	Alias string `json:"result"`
}

type PostHandler struct {
	BaseHandler
}

func NewPostHandler(app *app.App) *PostHandler {
	return &PostHandler{
		BaseHandler: BaseHandler{app},
	}
}

func (handler *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	URL, err := parseBody(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if alias, err := handler.app.Storage.GetAlias(ctx, storage.OriginalURL(URL)); err == nil {
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+string(alias), true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		alias := utils.RandomString(aliasSize)
		batch := map[storage.Alias]storage.OriginalURL{
			storage.Alias(alias): storage.OriginalURL(URL),
		}
		if err := handler.app.Storage.Add(ctx, batch); err != nil {
			log.Println("Can not add note to database")
			return
		}
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+alias, false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func parseBody(req *http.Request) (string, error) {
	if req.Header.Get("Content-Type") == "application/json" {
		jsonReq := new(postJSONRequest)
		err := json.NewDecoder(req.Body).Decode(jsonReq)
		if err != nil {
			return "", err
		}
		return jsonReq.URL, nil
	}
	body, err := io.ReadAll(req.Body)
	return string(body), err
}

func printResponse(w http.ResponseWriter, req *http.Request, alias string, allreadyAdded bool) error {
	if req.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		if allreadyAdded {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		err := json.NewEncoder(w).Encode(postJSONResponse{Alias: alias})
		return err
	}
	if allreadyAdded {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, alias)
	return nil
}
