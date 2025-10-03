// Package main реализует основное приложение для сокращения URL
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/cmd/shortener/router"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/services/ds"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
	"go.uber.org/zap"
)

const (
	// DeleteWorkers количество воркеров для удаления URL
	DeleteWorkers = 10

	// ShutdownTimeout таймаут для корректного завершения работы сервера
	ShutdownTimeout = 10 * time.Second

	// ServerTimeout таймаут для чтения и записи HTTP запросов
	ServerTimeout = 10 * time.Second
)

// main точка входа в приложение
func main() {
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	if err := run(logger.Sugar()); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}

// run основная функция запуска приложения
// logger - логгер для записи логов приложения
// Возвращает ошибку, если приложение не удалось запустить
func run(logger *zap.SugaredLogger) error {
	// Создание и парсинг конфигурации
	conf := config.NewConfig()
	if err := conf.ParseFlags(); err != nil {
		logger.Errorw("Failed to parse flags", "error", err)
		return err
	}

	// Создание хранилища
	store, err := storage.NewStorage(conf)
	if err != nil {
		logger.Errorw("Failed to create storage", "error", err)
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()
		if err := store.CloseStorage(ctx); err != nil {
			logger.Errorw("Failed to close storage", "error", err)
		}
	}()

	// Создание сервиса удаления
	deleteSvc := ds.NewDeleteService(store)
	deleteSvc.Start(DeleteWorkers)
	defer deleteSvc.Stop()

	// Создание приложения
	application := app.NewApp(store, conf, deleteSvc)

	// Создание маршрутизатора
	h := router.Build(application)

	// Создание HTTP сервера
	srv := &http.Server{
		Addr:         application.Config.ServerAddress,
		Handler:      h,
		ReadTimeout:  ServerTimeout,
		WriteTimeout: ServerTimeout,
	}

	// Каналы для обработки ошибок и сигналов
	errCh := make(chan error, 1)
	go func() {
		logger.Infow("Starting server", "address", application.Config.ServerAddress)
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание завершения
	select {
	case sig := <-sigCh:
		logger.Infow("Received signal, shutting down...", "signal", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Errorw("Server error", "error", err)
			return err
		}
	}

	// Корректное завершение работы
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorw("Server shutdown failed", "error", err)
		return err
	}

	logger.Info("Server shutdown completed")
	return nil
}
