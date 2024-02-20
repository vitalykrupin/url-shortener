package handlers

import (
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
)

type BaseHandler struct {
	app *config.App
}
