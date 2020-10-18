package main

import (
	"context"
	"garagesale/internal/platform/conf"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := conf.Parse()

	db, err := database.Open(config.Database)
	if err != nil {
		return errors.Wrap(err, "Connecting to db")
	}

	productsHandler := handlers.NewProductsHandler(db)

	api := http.Server{
		Addr:         config.Address(),
		Handler:      http.HandlerFunc(productsHandler.List),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "Starting server")

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout())
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", config.ShutdownTimeout(), err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "Server gracefully shutdown")
		}
	}

	return nil
}
