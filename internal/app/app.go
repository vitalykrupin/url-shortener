package app

import (
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/services/deleter"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type App struct {
	Config        *config.Config
	Storage       storage.Storage
	DeleteService *ds.DeleteService
}

func NewApp(config *config.Config, storage storage.Storage, deleteService *ds.DeleteService) *App {
	return &App{config, storage, deleteService}
}
