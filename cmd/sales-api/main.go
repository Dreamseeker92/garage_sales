package main

import (
	"context"
	"garagesale/internal/platform/conf"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/database"
)

func main() {
	config := conf.Parse()

	db, err := database.Open(config.Database)
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}

	productsHandler := handlers.Products{DB: db}

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
		log.Fatalf("error: starting server: %s", err)

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
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}
