package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app"
	"github.com/vitalykrupin/url-shortener.git/internal/app/handlers"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
)

func Route(app *app.App, conf *config.Config) {
	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.JwtAuthorization)

	router.Handle(`/{id}`, handlers.NewGetHandler(app))
	router.Handle(`/`, handlers.NewPostHandler(app))
	router.Method(http.MethodPost, `/api/shorten`, handlers.NewPostHandler(app))
	router.Method(http.MethodPost, `/api/shorten/batch`, handlers.NewPostBatchHandler(app))
	router.Method(http.MethodGet, `/ping`, handlers.NewGetPingHandler(app))
	router.Method(http.MethodGet, `/api/user/urls`, handlers.NewGetAllUserURLs(app))

	errListen := http.ListenAndServe(conf.ServerAddress, router)
	if errListen != nil {
		log.Fatalf("ListenAndServe returns error: %v", errListen)
	}
}
