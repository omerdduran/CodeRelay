package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"coderelay/backend/internal/server"
	"coderelay/backend/internal/storage"
	"coderelay/backend/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPath := envOrDefault("DB_PATH", "./coderelay.db")

	cfg := server.Config{
		Addr:   envOrDefault("API_ADDR", ":8080"),
		DBPath: dbPath,
	}

	// Initialize database for worker
	store, err := storage.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	// Run migrations
	if err := store.Migrate(); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	if err := store.Seed(); err != nil {
		// Seed errors are non-fatal (data might exist)
	}

	// Start worker in background
	w := worker.New(store)
	go w.Start(ctx)

	// Start HTTP server
	srv := server.New(cfg)
	log.Printf("starting api server on %s", cfg.Addr)

	if err := srv.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("api server exited: %v", err)
	}

	log.Println("api server shutdown complete")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
