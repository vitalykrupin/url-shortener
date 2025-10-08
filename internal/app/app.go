// Package app provides the main application structure and creation functions
package app

import (
	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener/internal/app/authservice"
	"github.com/vitalykrupin/url-shortener/internal/app/services/ds"
	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// App is the main application structure
type App struct {
	// Store is the interface for working with data storage
	Store storage.Storage

	// Config is the application configuration
	Config *config.Config

	// DeleteService is the service for deleting URLs
	DeleteService ds.DeleteServiceInterface

	// AuthService is the service for user authentication
	AuthService *authservice.AuthService
}

// NewApp creates a new application instance
// store is the interface for working with data storage
// conf is the application configuration
// deleteService is the service for deleting URLs
// authService is the service for user authentication
// Returns a pointer to App
func NewApp(store storage.Storage, conf *config.Config, deleteService ds.DeleteServiceInterface, authService *authservice.AuthService) *App {
	return &App{
		Store:         store,
		Config:        conf,
		DeleteService: deleteService,
		AuthService:   authService,
	}
}
