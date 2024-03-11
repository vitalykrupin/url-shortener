package main

import (
	"log"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/router"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/services/ds"
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
	store, err := storage.NewStorage(conf)
	if err != nil {
		log.Println("Can not create storage", err)
		return err
	}

	ds := ds.NewDeleteService(store)
	ds.Start()
	defer ds.Stop()

	router.Route(app.NewApp(store, conf, ds))

	return nil
}
