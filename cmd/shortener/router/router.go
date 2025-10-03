// Package router предоставляет функциональность для маршрутизации HTTP запросов
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/handlers"
	appMiddleware "github.com/vitalykrupin/url-shortener/internal/app/middleware"
)

// Build создает и настраивает маршрутизатор HTTP запросов
// app - экземпляр приложения
// Возвращает http.Handler для обработки запросов
func Build(app *app.App) http.Handler {
	r := chi.NewRouter()

	// Стандартные middleware из chi
	r.Use(chiMiddleware.RequestID)                 // Добавляет уникальный ID к каждому запросу
	r.Use(chiMiddleware.RealIP)                    // Определяет реальный IP адрес клиента
	r.Use(chiMiddleware.Logger)                    // Логирует HTTP запросы
	r.Use(chiMiddleware.Recoverer)                 // Восстанавливается после паники в обработчиках
	r.Use(chiMiddleware.Timeout(60 * time.Second)) // Устанавливает таймаут для запросов

	// Пользовательские middleware
	r.Use(appMiddleware.Logging)          // Пользовательское логирование
	r.Use(appMiddleware.GzipMiddleware)   // Поддержка сжатия gzip
	r.Use(appMiddleware.JwtAuthorization) // Авторизация через JWT

	// Маршруты для получения URL
	r.Handle(`/{id}`, handlers.NewGetHandler(app)) // Получение оригинального URL по короткому alias

	// Маршруты для проверки работоспособности
	r.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(app)) // Проверка подключения к базе данных

	// Маршруты для работы с URL пользователя
	r.Method(http.MethodGet, `/api/user/urls`, handlers.NewGetAllUserURLs(app)) // Получение всех URL пользователя

	// Маршруты для создания коротких URL
	r.Handle(`/`, handlers.NewPostHandler(app))                                        // Создание короткого URL из тела запроса
	r.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(app))            // Создание короткого URL из JSON
	r.Method(http.MethodPost, `/api/shorten/batch`, handlers.NewPostBatchHandler(app)) // Создание нескольких коротких URL

	// Маршруты для удаления URL
	r.Method(http.MethodDelete, `/api/user/urls`, handlers.NewDeleteHandler(app)) // Удаление URL пользователя

	return r
}
