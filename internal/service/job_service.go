package service

import (
	"context"
	"fmt"

	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
	"github.com/abhisheksainimitawa/job-aggregator/internal/repository"
	"github.com/abhisheksainimitawa/job-aggregator/internal/scraper"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

// JobService handles business logic for jobs
type JobService struct {
	repo    *repository.JobRepository
	scraper *scraper.Engine
}

// NewJobService creates a new job service
func NewJobService(repo *repository.JobRepository, scraperEngine *scraper.Engine) *JobService {
	return &JobService{
		repo:    repo,
		scraper: scraperEngine,
	}
}

// GetJob retrieves a job by ID
func (s *JobService) GetJob(ctx context.Context, id int64) (*models.Job, error) {
	return s.repo.FindByID(ctx, id)
}

// SearchJobs searches for jobs based on criteria
func (s *JobService) SearchJobs(ctx context.Context, query *models.JobSearchQuery) ([]*models.Job, error) {
	return s.repo.Search(ctx, query)
}

// GetStats retrieves job statistics
func (s *JobService) GetStats(ctx context.Context) (*models.JobStats, error) {
	return s.repo.GetStats(ctx)
}

// RunScraper runs the scraper and stores results
func (s *JobService) RunScraper(ctx context.Context, query string) (int, error) {
	logger.Info("Starting job scraper for query: %s", query)

	// Run the scraper
	jobs, err := s.scraper.Start(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("scraper failed: %w", err)
	}

	// Store jobs in database with deduplication
	if err := s.repo.CreateBatch(ctx, jobs); err != nil {
		return 0, fmt.Errorf("failed to store jobs: %w", err)
	}

	logger.Info("Successfully scraped and stored %d jobs", len(jobs))
	return len(jobs), nil
}

// GetScraperStats returns current scraper statistics
func (s *JobService) GetScraperStats() scraper.Stats {
	return s.scraper.GetStats()
}
