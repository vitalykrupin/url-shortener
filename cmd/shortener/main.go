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
)

const DeleteWorkers = 10

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf := config.NewConfig()
	conf.ParseFlags()

	var err error
	store, err := storage.NewStorage(conf)
	if err != nil {
		log.Println("Can not create storage", err)
		return err
	}

	deleteSvc := ds.NewDeleteService(store)
	deleteSvc.Start(DeleteWorkers)
	defer deleteSvc.Stop()

	application := app.NewApp(store, conf, deleteSvc)

	h := router.Build(application)

	srv := &http.Server{
		Addr:         application.Config.ServerAddress,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("Received signal: %v, shutting down...", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = store.CloseStorage(ctx)
	return srv.Shutdown(ctx)
}
