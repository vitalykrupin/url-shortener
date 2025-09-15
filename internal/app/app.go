package app

import (
	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app/services/ds"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

type App struct {
	Store         storage.Storage
	Config        *config.Config
	DeleteService ds.DeleteServiceInterface
}

func NewApp(store storage.Storage, conf *config.Config, deleteService ds.DeleteServiceInterface) *App {
	return &App{Store: store, Config: conf, DeleteService: deleteService}
}
