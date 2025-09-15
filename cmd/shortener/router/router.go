package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener/internal/app"
	"github.com/vitalykrupin/url-shortener/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener/internal/app/middleware"
)

func Build(app *app.App) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.JwtAuthorization)

	r.Handle(`/{id}`, handlers.NewGetHandler(app))
	r.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(app))
	r.Method(http.MethodGet, `/api/user/urls`, handlers.NewGetAllUserURLs(app))

	r.Handle(`/`, handlers.NewPostHandler(app))
	r.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(app))
	r.Method(http.MethodPost, `/api/shorten/batch`, handlers.NewPostBatchHandler(app))

	r.Method(http.MethodDelete, `/api/user/urls`, handlers.NewDeleteHandler(app))

	return r
}
