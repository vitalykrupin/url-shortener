package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func main() {
	store := storage.NewStorage()
	router := chi.NewRouter()

	router.Handle(`/`, handlers.NewPostHandler(store))
	router.Handle(`/{id}`, handlers.NewGetHandler(store))

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
