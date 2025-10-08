// Package handlers provides HTTP request handlers
package handlers

import (
	"time"

	"github.com/vitalykrupin/url-shortener/internal/app"
)

const (
	// ctxTimeout is the default context timeout
	ctxTimeout = 10 * time.Second
	
	// idParam is the URL parameter name for ID
	idParam = "id"
	
	// aliasSize is the default alias size
	aliasSize = 8
)

// BaseHandler is the base structure for all handlers
type BaseHandler struct {
	app *app.App
}

// NewBaseHandler is the constructor for BaseHandler
func NewBaseHandler(app *app.App) *BaseHandler {
	return &BaseHandler{app: app}
}
