package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func main() {
	config := config.Config{}
	config.ParseFlags()
	store := storage.NewStorage()
	router := chi.NewRouter()

	router.Handle(`/`, handlers.NewPostHandler(store, config))
	router.Handle(`/{id}`, handlers.NewGetHandler(store, config))

	err := http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		panic(err)
	}
}
