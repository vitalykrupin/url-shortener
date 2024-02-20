package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vitalykrupin/url-shortener.git/internal/app/storage"
)

type App struct {
	Config  *Config
	Storage storage.StorageKeeper
	DBPool  *pgxpool.Pool
}

func NewApp(config *Config, storage storage.StorageKeeper) *App {
	return &App{config, storage, nil}
}
