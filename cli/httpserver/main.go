package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"redirectware/internal"
	"redirectware/storage/postgres"
)

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		slog.Error("db open failed", "err", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		slog.Error("db ping failed", "err", err)
		os.Exit(1)
	}

	storage := postgres.New(db)

	// Create handler with storage
	handler := internal.NewHandler(storage)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("starting redirect service", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
