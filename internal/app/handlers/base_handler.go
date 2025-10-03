package handlers

import (
	"time"

	"github.com/vitalykrupin/url-shortener/internal/app"
)

const (
	ctxTimeout = 10 * time.Second
	idParam    = "id"
	aliasSize  = 8
)

// BaseHandler - базовая структура для всех обработчиков
type BaseHandler struct {
	app *app.App
}

// NewBaseHandler - конструктор для BaseHandler
func NewBaseHandler(app *app.App) *BaseHandler {
	return &BaseHandler{app: app}
}
