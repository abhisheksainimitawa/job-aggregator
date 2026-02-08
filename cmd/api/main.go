package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhisheksainimitawa/job-aggregator/internal/api"
	"github.com/abhisheksainimitawa/job-aggregator/internal/config"
	"github.com/abhisheksainimitawa/job-aggregator/internal/repository"
	"github.com/abhisheksainimitawa/job-aggregator/internal/scraper"
	"github.com/abhisheksainimitawa/job-aggregator/internal/service"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	logger.Info("Starting Job Aggregator API Server...")

	// Initialize database
	db, err := repository.NewDB(cfg.GetDatabaseDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := repository.InitSchema(db); err != nil {
		logger.Fatal("Failed to initialize database schema: %v", err)
	}

	// Initialize repositories
	jobRepo := repository.NewJobRepository(db)

	// Initialize scraper engine
	scraperEngine := scraper.NewEngine(cfg.Scraper.Workers, cfg.Scraper.RateLimit)
	scraperEngine.RegisterSource(scraper.NewIndeedScraper())
	scraperEngine.RegisterSource(scraper.NewLinkedInScraper())
	scraperEngine.RegisterSource(scraper.NewGlassdoorScraper())
	defer scraperEngine.Shutdown()

	// Initialize services
	jobService := service.NewJobService(jobRepo, scraperEngine)

	// Initialize HTTP handler
	handler := api.NewHandler(jobService)
	router := handler.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server listening on %s", cfg.GetServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped successfully")
}
