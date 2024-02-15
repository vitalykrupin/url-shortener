package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf := &config.Config{}
	conf.InitConfig()

	store, err := storage.NewDB(conf)
	if err != nil {
		store = storage.NewFileStorage(conf)
	}

	appInstance := app.NewApp(conf, store)

	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipMiddleware)

	router.Handle(`/{id}`, handlers.NewGetHandler(appInstance))
	router.Handle(`/`, handlers.NewPostHandler(appInstance))
	router.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(appInstance))
	router.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(appInstance))

	errListen := http.ListenAndServe(conf.ServerAddress, router)
	if errListen != nil {
		log.Fatalf("ListenAndServe returns error: %v", err)
	}
	return nil
}
