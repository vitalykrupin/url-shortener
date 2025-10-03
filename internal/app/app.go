// Package app предоставляет основную структуру приложения и функции для его создания
package app

import (
	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app/services/ds"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// App основная структура приложения
type App struct {
	// Store интерфейс для работы с хранилищем данных
	Store storage.Storage

	// Config конфигурация приложения
	Config *config.Config

	// DeleteService сервис для удаления URL
	DeleteService ds.DeleteServiceInterface
}

// NewApp создает новый экземпляр приложения
// store - интерфейс для работы с хранилищем данных
// conf - конфигурация приложения
// deleteService - сервис для удаления URL
// Возвращает указатель на App
func NewApp(store storage.Storage, conf *config.Config, deleteService ds.DeleteServiceInterface) *App {
	return &App{
		Store:         store,
		Config:        conf,
		DeleteService: deleteService,
	}
}
