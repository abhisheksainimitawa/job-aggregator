package models

import (
	"time"
)

// Job represents a job listing aggregated from various sources
type Job struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Company     string    `json:"company" db:"company"`
	Location    string    `json:"location" db:"location"`
	Salary      string    `json:"salary,omitempty" db:"salary"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Source      string    `json:"source" db:"source"` // indeed, linkedin, etc.
	RemoteOk    bool      `json:"remote_ok" db:"remote_ok"`
	JobType     string    `json:"job_type" db:"job_type"` // full-time, part-time, contract
	PostedAt    time.Time `json:"posted_at" db:"posted_at"`
	ScrapedAt   time.Time `json:"scraped_at" db:"scraped_at"`
	Hash        string    `json:"-" db:"hash"` // For deduplication
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// JobSearchQuery represents search parameters
type JobSearchQuery struct {
	Keywords  string
	Location  string
	Remote    *bool
	JobType   string
	Source    string
	MinSalary int
	Page      int
	Limit     int
}

// JobStats represents aggregated statistics
type JobStats struct {
	TotalJobs       int64            `json:"total_jobs"`
	JobsBySource    map[string]int64 `json:"jobs_by_source"`
	JobsByType      map[string]int64 `json:"jobs_by_type"`
	RemoteJobs      int64            `json:"remote_jobs"`
	TodayJobs       int64            `json:"today_jobs"`
	LastScrapedAt   time.Time        `json:"last_scraped_at"`
	TopCompanies    []CompanyCount   `json:"top_companies"`
	TopLocations    []LocationCount  `json:"top_locations"`
}

// CompanyCount represents job count by company
type CompanyCount struct {
	Company string `json:"company"`
	Count   int64  `json:"count"`
}

// LocationCount represents job count by location
type LocationCount struct {
	Location string `json:"location"`
	Count    int64  `json:"count"`
}

// ScraperStatus represents the current status of the scraper
type ScraperStatus struct {
	IsRunning      bool      `json:"is_running"`
	CurrentSource  string    `json:"current_source,omitempty"`
	JobsScraped    int       `json:"jobs_scraped"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	LastCompletedAt time.Time `json:"last_completed_at,omitempty"`
	ErrorCount     int       `json:"error_count"`
}
