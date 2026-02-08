package main

import (
	"context"
	"flag"
	"time"

	"github.com/abhisheksainimitawa/job-aggregator/internal/config"
	"github.com/abhisheksainimitawa/job-aggregator/internal/repository"
	"github.com/abhisheksainimitawa/job-aggregator/internal/scraper"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

func main() {
	// Parse command line flags
	query := flag.String("query", "golang developer", "Search query for jobs")
	source := flag.String("source", "", "Specific source to scrape (indeed, linkedin, glassdoor)")
	workers := flag.Int("workers", 10, "Number of concurrent workers")
	flag.Parse()

	logger.Info("Starting Job Scraper CLI...")
	logger.Info("Query: %s, Workers: %d", *query, *workers)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Override workers if specified
	if *workers > 0 {
		cfg.Scraper.Workers = *workers
	}

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

	// Initialize repository
	jobRepo := repository.NewJobRepository(db)

	// Initialize scraper engine
	scraperEngine := scraper.NewEngine(cfg.Scraper.Workers, cfg.Scraper.RateLimit)

	// Register sources based on filter
	if *source == "" || *source == "indeed" {
		scraperEngine.RegisterSource(scraper.NewIndeedScraper())
	}
	if *source == "" || *source == "linkedin" {
		scraperEngine.RegisterSource(scraper.NewLinkedInScraper())
	}
	if *source == "" || *source == "glassdoor" {
		scraperEngine.RegisterSource(scraper.NewGlassdoorScraper())
	}

	defer scraperEngine.Shutdown()

	// Run scraper with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger.Info("Starting scraping process...")
	jobs, err := scraperEngine.Start(ctx, *query)
	if err != nil {
		logger.Fatal("Scraping failed: %v", err)
	}

	logger.Info("Scraped %d jobs, storing in database...", len(jobs))

	// Store jobs in database with deduplication
	if err := jobRepo.CreateBatch(ctx, jobs); err != nil {
		logger.Fatal("Failed to store jobs: %v", err)
	}

	// Get statistics
	stats := scraperEngine.GetStats()
	logger.Info("=== Scraping Summary ===")
	logger.Info("Jobs Scraped: %d", stats.JobsScraped)
	logger.Info("Errors: %d", stats.Errors)
	logger.Info("Duration: %v", stats.EndTime.Sub(stats.StartTime))
	logger.Info("Jobs/Second: %.2f", float64(stats.JobsScraped)/stats.EndTime.Sub(stats.StartTime).Seconds())
	logger.Info("========================")

	logger.Info("Scraping completed successfully!")
}
