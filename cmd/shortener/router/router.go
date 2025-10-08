// Package router provides functionality for HTTP request routing
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

// Build creates and configures the HTTP request router
// app is the application instance
// Returns http.Handler for handling requests
func Build(app *app.App) http.Handler {
	r := chi.NewRouter()

	// Standard middleware from chi
	r.Use(chiMiddleware.RequestID)                 // Adds unique ID to each request
	r.Use(chiMiddleware.RealIP)                    // Determines real client IP address
	r.Use(chiMiddleware.Logger)                    // Logs HTTP requests
	r.Use(chiMiddleware.Recoverer)                 // Recovers from panics in handlers
	r.Use(chiMiddleware.Timeout(60 * time.Second)) // Sets timeout for requests

	// Custom middleware
	r.Use(appMiddleware.Logging)          // Custom logging
	r.Use(appMiddleware.GzipMiddleware)        // Gzip compression support
	r.Use(appMiddleware.ExternalJwtAuthorization) // External JWT authorization

	// Routes for getting URLs
	r.Handle(`/{id}`, handlers.NewGetHandler(app)) // Get original URL by short alias

	// Routes for health checks
	r.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(app)) // Check database connection

	// Routes for working with user URLs
	r.Method(http.MethodGet, `/api/user/urls`, handlers.NewGetAllUserURLs(app)) // Get all user URLs

	// Routes for creating short URLs
	r.Handle(`/`, handlers.NewPostHandler(app))                                        // Create short URL from request body
	r.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(app))            // Create short URL from JSON
	r.Method(http.MethodPost, `/api/shorten/batch`, handlers.NewPostBatchHandler(app)) // Create multiple short URLs

	// Routes for deleting URLs
	r.Method(http.MethodDelete, `/api/user/urls`, handlers.NewDeleteHandler(app)) // Delete user URLs

	return r
}
