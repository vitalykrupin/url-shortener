package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
	"github.com/vitalykrupin/url-shortener.git/internal/app/utils"
)

const (
	aliasSize int = 7
	idParam       = "id"
)

type postJSONRequest struct {
	FullURL string `json:"url"`
}

type postJSONResponse struct {
	Alias string `json:"result"`
}

type PostHandler struct {
	BaseHandler
}

func NewPostHandler(store *storage.Store, config config.Config) *PostHandler {
	return &PostHandler{
		BaseHandler: BaseHandler{
			store:  store,
			config: config,
		},
	}
}

func (handler *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fullURL, err := parseJSON(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if alias, ok := handler.store.Store.FullURLKeysMap[string(fullURL)]; ok {
		err := printResponse(w, req, handler.config.ResponseAddress+"/"+alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		alias := utils.RandomString(aliasSize)
		handler.store.Store.FullURLKeysMap[string(fullURL)] = alias
		handler.store.Store.AliasKeysMap[alias] = string(fullURL)
		err := printResponse(w, req, handler.config.ResponseAddress+"/"+alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	handler.store.SaveJSONtoFS(handler.config.FileStorePath)
}

func parseJSON(req *http.Request) (string, error) {
	if req.Header.Get("Content-Type") == "application/json" {
		jsonReq := new(postJSONRequest)
		err := json.NewDecoder(req.Body).Decode(jsonReq)
		if err != nil {
			return "", err
		}
		return jsonReq.FullURL, nil
	}
	body, err := io.ReadAll(req.Body)
	return string(body), err
}

func printResponse(w http.ResponseWriter, req *http.Request, alias string) error {
	if req.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(postJSONResponse{Alias: alias})
		return err
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, alias)
	return nil
}
