package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func main() {
	config := config.Config{}
	config.InitConfig()

	store := storage.NewStorage()
	storeLoadErr := store.LoadJSONfromFS(config.FileStorePath)
	if storeLoadErr != nil {
		log.Fatalf("Can not load store from file: %v", storeLoadErr)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipMiddleware)

	router.Handle(`/{id}`, handlers.NewGetHandler(store, config))
	router.Handle(`/`, handlers.NewPostHandler(store, config))
	router.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(store, config))

	err := http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		log.Fatalf("ListenAndServe returns error: %v", err)
	}
}
