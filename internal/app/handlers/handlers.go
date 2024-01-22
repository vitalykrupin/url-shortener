package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

const aliasSize int = 7

type PostHandler struct {
	http.Handler
	store *storage.DB
}

type GetHandler struct {
	http.Handler
	store *storage.DB
}

func NewPostHandler(store *storage.DB) *PostHandler {
	h := new(PostHandler)
	h.store = store
	return h
}

func NewGetHandler(store *storage.DB) *GetHandler {
	h := new(GetHandler)
	h.store = store
	return h
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
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
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
		fmt.Fprint(w, "http://"+req.Host+"/"+alias)
		return
	} else {
		alias := handler.randomString(aliasSize)
		handler.store.FullURLKeysMap[string(fullURL)] = alias
		handler.store.AliasKeysMap[alias] = string(fullURL)
		fmt.Fprint(w, "http://"+req.Host+"/"+handler.store.FullURLKeysMap[string(fullURL)])
	}
}

func (handler *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	alias, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "Get query require Id", http.StatusBadRequest)
		return
	}
	if fullURL, ok := handler.store.AliasKeysMap[alias]; ok {
		w.Header().Add("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Fprint(w, "URL by alias "+alias+" is not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
