package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reangeline/go-shipping-products/internal/app"
	"github.com/reangeline/go-shipping-products/internal/app/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Proccess to dependency injection
	container, err := app.Wire(cfg)
	if err != nil {
		log.Fatalf("wire failed: %v", err)
	}

	// I created a container.HTTP handler to remove depency of gin
	// If necessary to change in the future
	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      container.HTTP,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start
	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Print("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Print("bye")
}
