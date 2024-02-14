package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	
	store := storage.NewMemoryStorage()
	storeLoadErr := store.LoadJSONfromFS(conf.FileStorePath)
	if storeLoadErr != nil {
		log.Fatalf("Can not load store from file: %v", storeLoadErr)
	}

	appInstance := app.NewApp(conf, store)
	if conf.DBDSN != "" {
		conn, err := sql.Open("pgx", conf.DBDSN)
		if err != nil {
			return err
		}
		appInstance.DB = conn
	}

	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipMiddleware)

	router.Handle(`/{id}`, handlers.NewGetHandler(appInstance))
	router.Handle(`/`, handlers.NewPostHandler(appInstance))
	router.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(appInstance))
	router.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(appInstance))

	err := http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		log.Fatalf("ListenAndServe returns error: %v", err)
	}
	return nil
}
