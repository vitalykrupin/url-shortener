package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

func main() {
	conf := &config.Config{}
	conf.InitConfig()

	store := storage.NewMemoryStorage()
	storeLoadErr := store.LoadJSONfromFS(conf.FileStorePath)
	if storeLoadErr != nil {
		log.Fatalf("Can not load store from file: %v", storeLoadErr)
	}

	app := config.NewApp(conf, store)
	if conf.DBDSN != "" {
		conn, err := pgxpool.New(context.Background(), conf.DBDSN)
		if err != nil {
			log.Fatalf("Error connecting to PG: %v", err)
		}
		app.DBPool = conn

		defer conn.Close()
	}

	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipMiddleware)

	router.Handle(`/{id}`, handlers.NewGetHandler(app))
	router.Handle(`/`, handlers.NewPostHandler(app))
	router.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(app))
	router.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(app))

	err := http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		log.Fatalf("ListenAndServe returns error: %v", err)
	}
}
