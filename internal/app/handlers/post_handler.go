package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/utils"
)

const (
	aliasSize int = 7
	idParam       = "id"
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

func NewPostHandler(app *config.App) *PostHandler {
	return &PostHandler{
		BaseHandler: BaseHandler{
			app: app,
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
	URL, err := parseJSON(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if alias, ok := handler.app.Storage.GetAlias(URL); ok {
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		alias := utils.RandomString(aliasSize)
		handler.app.Storage.AddToMemoryStore(URL, alias)
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	handler.app.Storage.SaveJSONtoFS(handler.app.Config.FileStorePath)
}

func parseJSON(req *http.Request) (string, error) {
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
