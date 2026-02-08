package scraper

import (
	"context"
	"testing"
	"time"

	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
)

func TestEngine_Start(t *testing.T) {
	// Create a test engine
	engine := NewEngine(5, 100)

	// Register test scrapers
	engine.RegisterSource(NewIndeedScraper())
	engine.RegisterSource(NewLinkedInScraper())

	// Run the scraper
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobs, err := engine.Start(ctx, "golang developer")
	if err != nil {
		t.Fatalf("Engine.Start() error = %v", err)
	}

	if len(jobs) == 0 {
		t.Error("Expected jobs to be scraped, got 0")
	}

	// Verify stats
	stats := engine.GetStats()
	if stats.JobsScraped == 0 {
		t.Error("Expected JobsScraped > 0")
	}

	t.Logf("Scraped %d jobs with %d errors", stats.JobsScraped, stats.Errors)
}

func TestGenerateJobHash(t *testing.T) {
	job1 := &models.Job{
		Title:    "Go Developer",
		Company:  "Tech Corp",
		Location: "Remote",
	}

	job2 := &models.Job{
		Title:    "Go Developer",
		Company:  "Tech Corp",
		Location: "Remote",
	}

	job3 := &models.Job{
		Title:    "Python Developer",
		Company:  "Tech Corp",
		Location: "Remote",
	}

	hash1 := generateJobHash(job1)
	hash2 := generateJobHash(job2)
	hash3 := generateJobHash(job3)

	if hash1 != hash2 {
		t.Error("Expected same hash for identical jobs")
	}

	if hash1 == hash3 {
		t.Error("Expected different hash for different jobs")
	}
}

func TestIndeedScraper_Scrape(t *testing.T) {
	scraper := NewIndeedScraper()

	if scraper.Name() != "Indeed" {
		t.Errorf("Expected name 'Indeed', got '%s'", scraper.Name())
	}

	ctx := context.Background()
	jobs, err := scraper.Scrape(ctx, "golang")

	if err != nil {
		t.Fatalf("Scrape() error = %v", err)
	}

	if len(jobs) == 0 {
		t.Error("Expected jobs to be returned")
	}

	// Verify job structure
	for _, job := range jobs {
		if job.Title == "" {
			t.Error("Expected job to have a title")
		}
		if job.Company == "" {
			t.Error("Expected job to have a company")
		}
		if job.Source != "Indeed" {
			t.Errorf("Expected source 'Indeed', got '%s'", job.Source)
		}
	}
}
