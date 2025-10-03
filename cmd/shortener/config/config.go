// Package config предоставляет функциональность для работы с конфигурацией приложения
package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v10"
)

const (
	// defaultServerAddress адрес сервера по умолчанию
	defaultServerAddress = "localhost:8080"

	// defaultResponseAddress базовый URL для ответов по умолчанию
	defaultResponseAddress = "http://localhost:8080"

	// defaultDBDSN строка подключения к базе данных по умолчанию
	defaultDBDSN = ""
)

// Config структура для хранения конфигурации приложения
type Config struct {
	// ServerAddress адрес сервера
	ServerAddress string `env:"SERVER_ADDRESS"`

	// ResponseAddress базовый URL для ответов
	ResponseAddress string `env:"BASE_URL"`

	// FileStorePath путь к файлу хранилища
	FileStorePath string `env:"FILE_STORAGE_PATH"`

	// DBDSN строка подключения к базе данных
	DBDSN string `env:"DATABASE_DSN"`
}

// NewConfig создает новый экземпляр конфигурации с значениями по умолчанию
// Возвращает указатель на Config
func NewConfig() *Config {
	return &Config{
		ServerAddress:   defaultServerAddress,
		ResponseAddress: defaultResponseAddress,
		FileStorePath:   filepath.Join(os.TempDir(), "short-url-db.json"),
		DBDSN:           defaultDBDSN,
	}
}

// ParseFlags парсит флаги командной строки и переменные окружения
// Возвращает ошибку, если парсинг не удался или конфигурация невалидна
func (c *Config) ParseFlags() error {
	// Регистрация флагов командной строки
	flag.Func("a", "example: '-a localhost:8080'", func(addr string) error {
		c.ServerAddress = addr
		return nil
	})
	flag.Func("b", "example: '-b http://localhost:8000'", func(addr string) error {
		c.ResponseAddress = addr
		return nil
	})
	flag.Func("f", "example: '-f /tmp/testfile.json'", func(path string) error {
		c.FileStorePath = path
		return nil
	})
	flag.Func("d", "example: '-d postgres://postgres:pwd@localhost:5432/postgres?sslmode=disable'", func(dbAddr string) error {
		c.DBDSN = dbAddr
		return nil
	})
	flag.Parse()

	// Парсинг переменных окружения
	err := env.Parse(c)
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Валидация конфигурации
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

// Validate проверяет корректность конфигурации
// Возвращает ошибку, если конфигурация невалидна
func (c *Config) Validate() error {
	// Проверка адреса сервера
	if c.ServerAddress == "" {
		return fmt.Errorf("server address is required")
	}

	// Проверка базового URL
	if c.ResponseAddress == "" {
		return fmt.Errorf("response address is required")
	}

	// Проверка корректности URL
	if _, err := url.ParseRequestURI(c.ResponseAddress); err != nil {
		return fmt.Errorf("invalid response address format: %w", err)
	}

	// Проверка пути к файлу хранилища (если используется файловое хранилище)
	if c.DBDSN == "" && c.FileStorePath == "" {
		return fmt.Errorf("either database DSN or file storage path must be provided")
	}

	return nil
}
