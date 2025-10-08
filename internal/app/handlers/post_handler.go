// Package handlers provides HTTP request handlers
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
	"github.com/vitalykrupin/url-shortener/internal/app/utils"
)

// postJSONRequest represents the JSON request structure for POST handler
type postJSONRequest struct {
	URL string `json:"url"`
}

// postJSONResponse represents the JSON response structure for POST handler
type postJSONResponse struct {
	Alias string `json:"result"`
}

// PostHandler handles POST requests for URL creation
type PostHandler struct {
	BaseHandler
}

// NewPostHandler is the constructor for PostHandler
func NewPostHandler(app *app.App) *PostHandler {
	return &PostHandler{
		BaseHandler: BaseHandler{app},
	}
}

// ServeHTTP handles the HTTP request for URL creation
func (handler *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), ctxTimeout)
	defer cancel()

	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	URL, err := parseBody(req)
	if err != nil {
		log.Println("Can not parse body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Check if URL is empty
	if URL == "" {
		log.Println("Empty URL in request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if alias, err := handler.app.Store.GetAlias(ctx, storage.OriginalURL(URL)); err == nil {
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+string(alias), true)
		if err != nil {
			log.Println("Can not print response", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		alias := utils.RandomString(aliasSize)
		batch := map[storage.Alias]storage.OriginalURL{
			storage.Alias(alias): storage.OriginalURL(URL),
		}
		if err := handler.app.Store.Add(ctx, batch); err != nil {
			log.Println("Can not add note to database", err)
		}
		err := printResponse(w, req, handler.app.Config.ResponseAddress+"/"+alias, false)
		if err != nil {
			log.Println("Can not print response", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

// parseBody parses the request body
// req is the HTTP request
// Returns the URL string and an error if parsing failed
func parseBody(req *http.Request) (string, error) {
	defer req.Body.Close()
	if req.Header.Get("Content-Type") == "application/json" {
		jsonReq := new(postJSONRequest)
		err := json.NewDecoder(req.Body).Decode(jsonReq)
		if err != nil {
			return "", err
		}
		return jsonReq.URL, nil
	}
	body, err := io.ReadAll(req.Body)
	stringBody := string(body)
	if stringBody == "" {
		log.Println("No body in request")
		return "", fmt.Errorf("no body in request")
	}
	return stringBody, err
}

// printResponse prints the response
// w is the HTTP response writer
// req is the HTTP request
// alias is the short URL alias
// alreadyAdded indicates if the URL was already added
// Returns an error if printing failed
func printResponse(w http.ResponseWriter, req *http.Request, alias string, alreadyAdded bool) error {
	if req.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		if alreadyAdded {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		err := json.NewEncoder(w).Encode(postJSONResponse{Alias: alias})
		return err
	}
	w.Header().Set("Content-Type", "text/plain")
	if alreadyAdded {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	fmt.Fprint(w, alias)
	return nil
}
