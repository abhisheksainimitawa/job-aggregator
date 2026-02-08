package scraper

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
)

// IndeedScraper scrapes Indeed job board (mock implementation)
type IndeedScraper struct {
	baseURL string
}

// NewIndeedScraper creates a new Indeed scraper
func NewIndeedScraper() *IndeedScraper {
	return &IndeedScraper{
		baseURL: "https://www.indeed.com",
	}
}

// Name returns the scraper name
func (s *IndeedScraper) Name() string {
	return "Indeed"
}

// Scrape scrapes jobs from Indeed (mock implementation for demo)
func (s *IndeedScraper) Scrape(ctx context.Context, query string) ([]*models.Job, error) {
	// In production, this would make actual HTTP requests to Indeed
	// For now, we'll generate mock data to demonstrate the architecture

	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	jobs := make([]*models.Job, 0)

	// Generate mock jobs
	jobTitles := []string{
		"Senior Go Developer",
		"Backend Engineer - Golang",
		"Full Stack Developer (Go/React)",
		"DevOps Engineer",
		"Site Reliability Engineer",
		"Cloud Platform Engineer",
		"Microservices Developer",
		"Software Engineer - Backend",
	}

	companies := []string{
		"Tech Corp", "StartupXYZ", "CloudSystems Inc",
		"DataFlow Technologies", "WebScale Solutions",
	}

	locations := []string{
		"San Francisco, CA", "New York, NY", "Remote",
		"Austin, TX", "Seattle, WA", "Boston, MA",
	}

	jobTypes := []string{"Full-time", "Contract", "Part-time"}

	// Generate 5-10 jobs
	numJobs := 5 + rand.Intn(6)
	for i := 0; i < numJobs; i++ {
		job := &models.Job{
			Title:       jobTitles[rand.Intn(len(jobTitles))],
			Company:     companies[rand.Intn(len(companies))],
			Location:    locations[rand.Intn(len(locations))],
			Description: fmt.Sprintf("Looking for an experienced developer with %d+ years in Go. %s", rand.Intn(5)+1, query),
			URL:         fmt.Sprintf("https://indeed.com/job/%d", rand.Intn(100000)),
			Source:      "Indeed",
			RemoteOk:    rand.Float32() > 0.5,
			JobType:     jobTypes[rand.Intn(len(jobTypes))],
			PostedAt:    time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
		}

		if rand.Float32() > 0.3 {
			job.Salary = fmt.Sprintf("$%dk - $%dk", 100+rand.Intn(100), 150+rand.Intn(100))
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// LinkedInScraper scrapes LinkedIn job board (mock implementation)
type LinkedInScraper struct {
	baseURL string
}

// NewLinkedInScraper creates a new LinkedIn scraper
func NewLinkedInScraper() *LinkedInScraper {
	return &LinkedInScraper{
		baseURL: "https://www.linkedin.com",
	}
}

// Name returns the scraper name
func (s *LinkedInScraper) Name() string {
	return "LinkedIn"
}

// Scrape scrapes jobs from LinkedIn (mock implementation for demo)
func (s *LinkedInScraper) Scrape(ctx context.Context, query string) ([]*models.Job, error) {
	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(700)) * time.Millisecond)

	jobs := make([]*models.Job, 0)

	jobTitles := []string{
		"Golang Software Engineer",
		"Backend Developer - Go",
		"Principal Engineer",
		"Staff Software Engineer",
		"Distributed Systems Engineer",
		"API Platform Engineer",
	}

	companies := []string{
		"Meta", "Google", "Amazon", "Microsoft",
		"Netflix", "Uber", "Airbnb",
	}

	locations := []string{
		"Menlo Park, CA", "Mountain View, CA", "Remote (US)",
		"Chicago, IL", "Denver, CO",
	}

	// Generate 3-8 jobs
	numJobs := 3 + rand.Intn(6)
	for i := 0; i < numJobs; i++ {
		job := &models.Job{
			Title:       jobTitles[rand.Intn(len(jobTitles))],
			Company:     companies[rand.Intn(len(companies))],
			Location:    locations[rand.Intn(len(locations))],
			Description: fmt.Sprintf("We are looking for talented engineers with expertise in %s", query),
			URL:         fmt.Sprintf("https://linkedin.com/jobs/view/%d", rand.Intn(1000000)),
			Source:      "LinkedIn",
			RemoteOk:    rand.Float32() > 0.4,
			JobType:     "Full-time",
			PostedAt:    time.Now().Add(-time.Duration(rand.Intn(14)) * 24 * time.Hour),
		}

		if rand.Float32() > 0.2 {
			job.Salary = fmt.Sprintf("$%dk - $%dk", 120+rand.Intn(150), 200+rand.Intn(150))
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GlassdoorScraper scrapes Glassdoor (mock implementation)
type GlassdoorScraper struct {
	baseURL string
}

// NewGlassdoorScraper creates a new Glassdoor scraper
func NewGlassdoorScraper() *GlassdoorScraper {
	return &GlassdoorScraper{
		baseURL: "https://www.glassdoor.com",
	}
}

// Name returns the scraper name
func (s *GlassdoorScraper) Name() string {
	return "Glassdoor"
}

// Scrape scrapes jobs from Glassdoor (mock implementation)
func (s *GlassdoorScraper) Scrape(ctx context.Context, query string) ([]*models.Job, error) {
	time.Sleep(time.Duration(rand.Intn(600)) * time.Millisecond)

	jobs := make([]*models.Job, 0)

	jobTitles := []string{
		"Senior Backend Engineer",
		"Go Developer",
		"Infrastructure Engineer",
		"Platform Engineer",
	}

	companies := []string{
		"Stripe", "Square", "Coinbase", "Robinhood",
		"Shopify", "Twilio",
	}

	// Generate 4-7 jobs
	numJobs := 4 + rand.Intn(4)
	for i := 0; i < numJobs; i++ {
		job := &models.Job{
			Title:       jobTitles[rand.Intn(len(jobTitles))],
			Company:     companies[rand.Intn(len(companies))],
			Location:    "Remote",
			Description: fmt.Sprintf("Join our team working with %s and cutting-edge technology", query),
			URL:         fmt.Sprintf("https://glassdoor.com/job-listing/%d", rand.Intn(500000)),
			Source:      "Glassdoor",
			RemoteOk:    true,
			JobType:     "Full-time",
			PostedAt:    time.Now().Add(-time.Duration(rand.Intn(20)) * 24 * time.Hour),
			Salary:      fmt.Sprintf("$%dk - $%dk", 110+rand.Intn(120), 180+rand.Intn(120)),
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
