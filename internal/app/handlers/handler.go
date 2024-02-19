package handlers

import (
	"time"

	"github.com/vitalykrupin/url-shortener.git/internal/app"
)

const (
	aliasSize  int = 7
	idParam        = "id"
	ctxTimeout     = 5 * time.Second
)

type BaseHandler struct {
	app *app.App
}
