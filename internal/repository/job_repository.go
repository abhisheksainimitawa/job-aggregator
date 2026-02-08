package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

// JobRepository handles database operations for jobs
type JobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create inserts a new job into the database
func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (title, company, location, salary, description, url, source, 
		                  remote_ok, job_type, posted_at, scraped_at, hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		job.Title, job.Company, job.Location, job.Salary, job.Description,
		job.URL, job.Source, job.RemoteOk, job.JobType, job.PostedAt,
		job.ScrapedAt, job.Hash, now, now,
	).Scan(&job.ID)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

// CreateBatch inserts multiple jobs in a single transaction (for performance)
func (r *JobRepository) CreateBatch(ctx context.Context, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO jobs (title, company, location, salary, description, url, source, 
		                  remote_ok, job_type, posted_at, scraped_at, hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (hash) DO UPDATE SET
			updated_at = EXCLUDED.updated_at,
			scraped_at = EXCLUDED.scraped_at
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, job := range jobs {
		err := stmt.QueryRowContext(ctx,
			job.Title, job.Company, job.Location, job.Salary, job.Description,
			job.URL, job.Source, job.RemoteOk, job.JobType, job.PostedAt,
			job.ScrapedAt, job.Hash, now, now,
		).Scan(&job.ID)

		if err != nil {
			logger.Error("Failed to insert job: %v", err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByID retrieves a job by ID
func (r *JobRepository) FindByID(ctx context.Context, id int64) (*models.Job, error) {
	query := `
		SELECT id, title, company, location, salary, description, url, source,
		       remote_ok, job_type, posted_at, scraped_at, hash, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`

	job := &models.Job{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Title, &job.Company, &job.Location, &job.Salary,
		&job.Description, &job.URL, &job.Source, &job.RemoteOk, &job.JobType,
		&job.PostedAt, &job.ScrapedAt, &job.Hash, &job.CreatedAt, &job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}

	return job, nil
}

// Search searches for jobs based on query parameters
func (r *JobRepository) Search(ctx context.Context, query *models.JobSearchQuery) ([]*models.Job, error) {
	sql := `
		SELECT id, title, company, location, salary, description, url, source,
		       remote_ok, job_type, posted_at, scraped_at, created_at, updated_at
		FROM jobs
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if query.Keywords != "" {
		sql += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d OR company ILIKE $%d)", argPos, argPos, argPos)
		args = append(args, "%"+query.Keywords+"%")
		argPos++
	}

	if query.Location != "" {
		sql += fmt.Sprintf(" AND location ILIKE $%d", argPos)
		args = append(args, "%"+query.Location+"%")
		argPos++
	}

	if query.Remote != nil {
		sql += fmt.Sprintf(" AND remote_ok = $%d", argPos)
		args = append(args, *query.Remote)
		argPos++
	}

	if query.JobType != "" {
		sql += fmt.Sprintf(" AND job_type = $%d", argPos)
		args = append(args, query.JobType)
		argPos++
	}

	if query.Source != "" {
		sql += fmt.Sprintf(" AND source = $%d", argPos)
		args = append(args, query.Source)
		argPos++
	}

	sql += " ORDER BY posted_at DESC"

	// Pagination
	if query.Limit == 0 {
		query.Limit = 20
	}
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, query.Limit, query.Page*query.Limit)

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]*models.Job, 0)
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID, &job.Title, &job.Company, &job.Location, &job.Salary,
			&job.Description, &job.URL, &job.Source, &job.RemoteOk, &job.JobType,
			&job.PostedAt, &job.ScrapedAt, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			logger.Error("Failed to scan job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetStats retrieves aggregated job statistics
func (r *JobRepository) GetStats(ctx context.Context) (*models.JobStats, error) {
	stats := &models.JobStats{
		JobsBySource: make(map[string]int64),
		JobsByType:   make(map[string]int64),
		TopCompanies: make([]models.CompanyCount, 0),
		TopLocations: make([]models.LocationCount, 0),
	}

	// Total jobs
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs").Scan(&stats.TotalJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to get total jobs: %w", err)
	}

	// Jobs by source
	rows, err := r.db.QueryContext(ctx, "SELECT source, COUNT(*) FROM jobs GROUP BY source")
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by source: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		var count int64
		if err := rows.Scan(&source, &count); err == nil {
			stats.JobsBySource[source] = count
		}
	}

	// Jobs by type
	rows, err = r.db.QueryContext(ctx, "SELECT job_type, COUNT(*) FROM jobs GROUP BY job_type")
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var jobType string
		var count int64
		if err := rows.Scan(&jobType, &count); err == nil {
			stats.JobsByType[jobType] = count
		}
	}

	// Remote jobs
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs WHERE remote_ok = true").Scan(&stats.RemoteJobs)
	if err != nil {
		logger.Error("Failed to get remote jobs count: %v", err)
	}

	// Today's jobs
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs WHERE DATE(posted_at) = CURRENT_DATE").Scan(&stats.TodayJobs)
	if err != nil {
		logger.Error("Failed to get today's jobs count: %v", err)
	}

	// Last scraped time
	err = r.db.QueryRowContext(ctx, "SELECT MAX(scraped_at) FROM jobs").Scan(&stats.LastScrapedAt)
	if err != nil {
		logger.Error("Failed to get last scraped time: %v", err)
	}

	// Top companies
	rows, err = r.db.QueryContext(ctx, `
		SELECT company, COUNT(*) as cnt 
		FROM jobs 
		GROUP BY company 
		ORDER BY cnt DESC 
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cc models.CompanyCount
			if err := rows.Scan(&cc.Company, &cc.Count); err == nil {
				stats.TopCompanies = append(stats.TopCompanies, cc)
			}
		}
	}

	// Top locations
	rows, err = r.db.QueryContext(ctx, `
		SELECT location, COUNT(*) as cnt 
		FROM jobs 
		GROUP BY location 
		ORDER BY cnt DESC 
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var lc models.LocationCount
			if err := rows.Scan(&lc.Location, &lc.Count); err == nil {
				stats.TopLocations = append(stats.TopLocations, lc)
			}
		}
	}

	return stats, nil
}

// DeleteOld deletes jobs older than the specified number of days
func (r *JobRepository) DeleteOld(ctx context.Context, days int) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM jobs WHERE posted_at < NOW() - INTERVAL '$1 days'",
		days,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old jobs: %w", err)
	}

	count, _ := result.RowsAffected()
	return count, nil
}

// ExistsByHash checks if a job with the given hash already exists
func (r *JobRepository) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM jobs WHERE hash = $1)",
		hash,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check job existence: %w", err)
	}

	return exists, nil
}
