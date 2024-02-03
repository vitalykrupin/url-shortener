package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

const (
	aliasSize int = 7
	idParam       = "id"
)

type PostHandler struct {
	http.Handler
	store  *storage.DB
	config config.Config
}

type GetHandler struct {
	http.Handler
	store  *storage.DB
	config config.Config
}

func NewPostHandler(store *storage.DB, config config.Config) *PostHandler {
	return &PostHandler{
		store:  store,
		config: config,
	}
}

func NewGetHandler(store *storage.DB, config config.Config) *GetHandler {
	return &GetHandler{
		store:  store,
		config: config,
	}
}

func (handler *PostHandler) randomString(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]rune, size)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

func (handler *PostHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Println("Only POST requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fullURL, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	if alias, ok := handler.store.FullURLKeysMap[string(fullURL)]; ok {
		fmt.Fprint(w, handler.config.ResponseAddress+"/"+alias)
		return
	} else {
		alias := handler.randomString(aliasSize)
		handler.store.FullURLKeysMap[string(fullURL)] = alias
		handler.store.AliasKeysMap[alias] = string(fullURL)
		fmt.Fprint(w, handler.config.ResponseAddress+"/"+alias)
	}
}

func (handler *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("Only GET requests are allowed!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	alias := chi.URLParam(req, idParam)
	if alias == "" {
		log.Println("Get query require Id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if fullURL, ok := handler.store.AliasKeysMap[alias]; ok {
		w.Header().Add("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		log.Println("URL by alias " + alias + " is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
