package app

import (
	"database/sql"

	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type App struct {
	Config  *config.Config
	Storage storage.StorageInterface
	DB      *sql.DB
}

func NewApp(config *config.Config, storage storage.StorageInterface) *App {
	return &App{config, storage, nil}
}
