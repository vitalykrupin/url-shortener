package app

import (
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type App struct {
	Config  *config.Config
	Storage storage.Storage
}

func NewApp(config *config.Config, storage storage.Storage) *App {
	return &App{config, storage}
}
