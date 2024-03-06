package app

import (
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type App struct {
	Store  storage.Storage
	Config *config.Config
}

func NewApp(store storage.Storage, conf *config.Config) *App {
	return &App{Store: store, Config: conf}
}
