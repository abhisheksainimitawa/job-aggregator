package scraper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/ratelimit"
)

// JobSource defines the interface for job board scrapers
type JobSource interface {
	Name() string
	Scrape(ctx context.Context, query string) ([]*models.Job, error)
}

// Engine is the main scraper engine with worker pools
type Engine struct {
	sources     []JobSource
	rateLimiter *ratelimit.RateLimiter
	workers     int
	jobCh       chan *models.Job
	errCh       chan error
	wg          sync.WaitGroup
	mu          sync.Mutex
	stats       Stats
}

// Stats holds scraping statistics
type Stats struct {
	JobsScraped int
	Errors      int
	StartTime   time.Time
	EndTime     time.Time
}

// NewEngine creates a new scraper engine
func NewEngine(workers int, rateLimit int) *Engine {
	return &Engine{
		sources:     make([]JobSource, 0),
		rateLimiter: ratelimit.NewRateLimiter(rateLimit),
		workers:     workers,
		jobCh:       make(chan *models.Job, 100),
		errCh:       make(chan error, 100),
		stats:       Stats{},
	}
}

// RegisterSource registers a job source scraper
func (e *Engine) RegisterSource(source JobSource) {
	e.sources = append(e.sources, source)
	logger.Info("Registered scraper source: %s", source.Name())
}

// Start starts the scraping engine with concurrent workers
func (e *Engine) Start(ctx context.Context, query string) ([]*models.Job, error) {
	e.stats.StartTime = time.Now()
	logger.Info("Starting scraper engine with %d workers for query: %s", e.workers, query)

	// Create worker pool
	jobsCh := make(chan *models.Job, 1000)
	done := make(chan struct{})

	// Collector goroutine
	var jobs []*models.Job
	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for job := range jobsCh {
			jobs = append(jobs, job)
			e.incrementJobCount()
		}
	}()

	// Error collector
	var errors []error
	go func() {
		for err := range e.errCh {
			errors = append(errors, err)
			e.incrementErrorCount()
			logger.Error("Scraper error: %v", err)
		}
	}()

	// Scrape each source concurrently
	sourceCh := make(chan JobSource, len(e.sources))
	for _, source := range e.sources {
		sourceCh <- source
	}
	close(sourceCh)

	// Start workers
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.worker(ctx, i, sourceCh, jobsCh, query)
	}

	// Wait for all workers to finish
	go func() {
		e.wg.Wait()
		close(jobsCh)
		close(e.errCh)
		close(done)
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		collectorWg.Wait()
		e.stats.EndTime = time.Now()
		logger.Info("Scraping completed: %d jobs, %d errors, duration: %v",
			e.stats.JobsScraped, e.stats.Errors, e.stats.EndTime.Sub(e.stats.StartTime))
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return jobs, nil
}

// worker is a worker goroutine that processes job sources
func (e *Engine) worker(ctx context.Context, id int, sources <-chan JobSource, jobs chan<- *models.Job, query string) {
	defer e.wg.Done()

	for source := range sources {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Info("Worker %d: Scraping %s", id, source.Name())

			// Rate limiting
			if err := e.rateLimiter.Wait(ctx); err != nil {
				e.errCh <- fmt.Errorf("rate limiter error: %w", err)
				continue
			}

			// Scrape the source
			sourceJobs, err := source.Scrape(ctx, query)
			if err != nil {
				e.errCh <- fmt.Errorf("%s scraper failed: %w", source.Name(), err)
				continue
			}

			// Send jobs to collector
			for _, job := range sourceJobs {
				// Add hash for deduplication
				job.Hash = generateJobHash(job)
				job.ScrapedAt = time.Now()

				select {
				case jobs <- job:
				case <-ctx.Done():
					return
				}
			}

			logger.Info("Worker %d: Scraped %d jobs from %s", id, len(sourceJobs), source.Name())
		}
	}
}

// generateJobHash creates a unique hash for job deduplication
func generateJobHash(job *models.Job) string {
	data := fmt.Sprintf("%s|%s|%s", job.Title, job.Company, job.Location)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// incrementJobCount increments the job counter
func (e *Engine) incrementJobCount() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stats.JobsScraped++
}

// incrementErrorCount increments the error counter
func (e *Engine) incrementErrorCount() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stats.Errors++
}

// GetStats returns current scraping statistics
func (e *Engine) GetStats() Stats {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.stats
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown() {
	logger.Info("Shutting down scraper engine")
	e.rateLimiter.Stop()
}
