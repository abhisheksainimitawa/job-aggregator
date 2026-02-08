package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

// NewDB creates a new database connection
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return db, nil
}

// InitSchema initializes the database schema
func InitSchema(db *sql.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS jobs (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			company VARCHAR(255) NOT NULL,
			location VARCHAR(255) NOT NULL,
			salary VARCHAR(100),
			description TEXT NOT NULL,
			url TEXT NOT NULL,
			source VARCHAR(50) NOT NULL,
			remote_ok BOOLEAN DEFAULT FALSE,
			job_type VARCHAR(50) NOT NULL,
			posted_at TIMESTAMP NOT NULL,
			scraped_at TIMESTAMP NOT NULL,
			hash VARCHAR(64) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_jobs_title ON jobs(title);
		CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);
		CREATE INDEX IF NOT EXISTS idx_jobs_location ON jobs(location);
		CREATE INDEX IF NOT EXISTS idx_jobs_source ON jobs(source);
		CREATE INDEX IF NOT EXISTS idx_jobs_posted_at ON jobs(posted_at DESC);
		CREATE INDEX IF NOT EXISTS idx_jobs_remote_ok ON jobs(remote_ok);
		CREATE INDEX IF NOT EXISTS idx_jobs_job_type ON jobs(job_type);
		CREATE INDEX IF NOT EXISTS idx_jobs_hash ON jobs(hash);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.Info("Database schema initialized successfully")
	return nil
}
