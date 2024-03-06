package main

import (
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
	conf := config.NewConfig()
	conf.ParseFlags()

	var err error
	storage.Store, err = storage.NewStorage(conf)
	if err != nil {
		log.Println("Can not create storage", err)
		return err
	}

	ds.DelService = ds.NewDeleteService()

	appInstance := app.NewApp(storage.Store, conf)

	router.Route(appInstance)

	return nil
}
