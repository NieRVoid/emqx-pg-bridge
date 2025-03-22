package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/NieRVoid/emqx-pg-bridge/internal/config"
	"github.com/NieRVoid/emqx-pg-bridge/internal/database"
	"github.com/NieRVoid/emqx-pg-bridge/internal/handler"
	"github.com/NieRVoid/emqx-pg-bridge/internal/processor"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

func main() {
	// Define command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configPath != "" {
		// Load from specified config file
		cfg, err = config.LoadFromFile(*configPath)
	} else {
		// Try to load from default locations
		cfg, err = config.Load()
	}

	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Starting EMQX-PostgreSQL Bridge",
		"version", cfg.Meta.Version,
		"buildDate", cfg.Meta.BuildDate)

	// Create a context for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to PostgreSQL
	db, err := database.NewPostgres(ctx, cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize processor registry
	registry := processor.NewProcessorRegistry(log)

	// Initialize and register processors
	centerProcessor := processor.NewCenterProcessor(db.Pool, log)
	normalProcessor := processor.NewNormalProcessor(db.Pool, log)

	registry.Register(centerProcessor)
	registry.Register(normalProcessor)

	// Setup HTTP router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Create webhook handler
	webhookHandler := handler.NewWebhookHandler(registry, log)

	// Register routes
	r.Post("/webhook", webhookHandler.Handle)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"%s","time":"%s"}`,
			cfg.Meta.Version, time.Now().Format(time.RFC3339))
	})

	// Start HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		IdleTimeout:  cfg.GetIdleTimeout(),
	}

	// Run server in a goroutine so that it doesn't block shutdown handling
	go func() {
		log.Info("Starting HTTP server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create shutdown context with 10 second timeout
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited properly")
}
