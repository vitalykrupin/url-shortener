package main

import (
	"context"
	"log"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/router"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/services/deleter"
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

	store, err := storage.NewDB(context.Background(), conf)
	if err != nil {
		store = storage.NewFileStorage(conf)
	}

	deleteService := ds.NewDeleteService()

	appInstance := app.NewApp(conf, store, deleteService)

	router.Route(appInstance, conf)

	return nil
}
