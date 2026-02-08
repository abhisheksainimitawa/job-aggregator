# API Examples

This document provides examples of how to use the Job Aggregator API.

## Prerequisites

Make sure the API server is running:
```bash
go run cmd/api/main.go
```

Or with Docker:
```bash
docker-compose up
```

## API Endpoints

### 1. Health Check

Check if the API is running:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "service": "job-aggregator",
  "time": "2026-02-09T10:30:00Z"
}
```

### 2. Run Scraper

Trigger the scraper to collect jobs:

```bash
curl -X POST http://localhost:8080/api/v1/scraper/run \
  -H "Content-Type: application/json" \
  -d '{"query": "golang developer"}'
```

Response:
```json
{
  "message": "Scraper completed successfully",
  "jobs_scraped": 25,
  "query": "golang developer"
}
```

### 3. List All Jobs

Get paginated list of jobs:

```bash
curl "http://localhost:8080/api/v1/jobs?page=0&limit=10"
```

Response:
```json
{
  "jobs": [
    {
      "id": 1,
      "title": "Senior Go Developer",
      "company": "Tech Corp",
      "location": "San Francisco, CA",
      "salary": "$150k - $200k",
      "description": "Looking for an experienced developer...",
      "url": "https://indeed.com/job/12345",
      "source": "Indeed",
      "remote_ok": true,
      "job_type": "Full-time",
      "posted_at": "2026-02-08T15:30:00Z",
      "scraped_at": "2026-02-09T10:00:00Z"
    }
  ],
  "page": 0,
  "limit": 10,
  "total": 10
}
```

### 4. Search Jobs

Search for specific jobs:

```bash
# Search by keywords
curl "http://localhost:8080/api/v1/jobs/search?q=golang&limit=5"

# Search with location
curl "http://localhost:8080/api/v1/jobs/search?q=backend&location=remote"

# Search for remote jobs only
curl "http://localhost:8080/api/v1/jobs/search?q=developer&remote=true"

# Search by job type
curl "http://localhost:8080/api/v1/jobs/search?type=Full-time&source=LinkedIn"
```

Response:
```json
{
  "jobs": [...],
  "query": "golang",
  "total": 15,
  "page": 0,
  "limit": 20
}
```

### 5. Get Single Job

Get details of a specific job:

```bash
curl "http://localhost:8080/api/v1/jobs/123"
```

Response:
```json
{
  "id": 123,
  "title": "Backend Engineer - Go",
  "company": "StartupXYZ",
  "location": "Remote",
  "salary": "$120k - $180k",
  "description": "We are looking for talented engineers...",
  "url": "https://linkedin.com/jobs/view/456789",
  "source": "LinkedIn",
  "remote_ok": true,
  "job_type": "Full-time",
  "posted_at": "2026-02-07T09:00:00Z",
  "scraped_at": "2026-02-09T10:00:00Z"
}
```

### 6. Get Statistics

Get aggregated job statistics:

```bash
curl "http://localhost:8080/api/v1/jobs/stats"
```

Response:
```json
{
  "total_jobs": 1543,
  "jobs_by_source": {
    "Indeed": 623,
    "LinkedIn": 512,
    "Glassdoor": 408
  },
  "jobs_by_type": {
    "Full-time": 1234,
    "Contract": 245,
    "Part-time": 64
  },
  "remote_jobs": 892,
  "today_jobs": 127,
  "last_scraped_at": "2026-02-09T10:15:00Z",
  "top_companies": [
    {"company": "Google", "count": 45},
    {"company": "Meta", "count": 38},
    {"company": "Amazon", "count": 35}
  ],
  "top_locations": [
    {"location": "Remote", "count": 892},
    {"location": "San Francisco, CA", "count": 245},
    {"location": "New York, NY", "count": 198}
  ]
}
```

### 7. Get Scraper Status

Check current scraper status:

```bash
curl "http://localhost:8080/api/v1/scraper/status"
```

Response:
```json
{
  "jobs_scraped": 25,
  "errors": 0,
  "start_time": "2026-02-09T10:00:00Z",
  "end_time": "2026-02-09T10:00:45Z"
}
```

## CLI Examples

### Run Scraper from Command Line

```bash
# Basic usage
go run cmd/scraper/main.go

# Custom query
go run cmd/scraper/main.go -query "python developer"

# Specific source only
go run cmd/scraper/main.go -source indeed -query "devops engineer"

# Custom worker count
go run cmd/scraper/main.go -workers 20 -query "full stack developer"
```

Output:
```
[2026-02-09 10:00:00] INFO: Starting Job Scraper CLI...
[2026-02-09 10:00:00] INFO: Query: golang developer, Workers: 10
[2026-02-09 10:00:00] INFO: Registered scraper source: Indeed
[2026-02-09 10:00:00] INFO: Registered scraper source: LinkedIn
[2026-02-09 10:00:00] INFO: Registered scraper source: Glassdoor
[2026-02-09 10:00:01] INFO: Starting scraper engine with 10 workers for query: golang developer
[2026-02-09 10:00:05] INFO: Scraping completed: 25 jobs, 0 errors, duration: 4.2s
[2026-02-09 10:00:05] INFO: === Scraping Summary ===
[2026-02-09 10:00:05] INFO: Jobs Scraped: 25
[2026-02-09 10:00:05] INFO: Errors: 0
[2026-02-09 10:00:05] INFO: Duration: 4.2s
[2026-02-09 10:00:05] INFO: Jobs/Second: 5.95
[2026-02-09 10:00:05] INFO: ========================
```

## Advanced Usage

### Combining Multiple Filters

```bash
curl "http://localhost:8080/api/v1/jobs/search?q=senior%20engineer&location=san%20francisco&remote=false&type=Full-time&page=0&limit=10"
```

### Using with jq for JSON Processing

```bash
# Get only job titles
curl -s "http://localhost:8080/api/v1/jobs?limit=5" | jq '.jobs[].title'

# Count remote jobs
curl -s "http://localhost:8080/api/v1/jobs/stats" | jq '.remote_jobs'

# List all companies hiring
curl -s "http://localhost:8080/api/v1/jobs/stats" | jq '.top_companies[].company'
```

### Automating Scraper with Cron

Add to your crontab to scrape every 6 hours:
```bash
0 */6 * * * cd /path/to/job-aggregator && /usr/local/go/bin/go run cmd/scraper/main.go -query "golang developer" >> /var/log/scraper.log 2>&1
```

## Testing the API

Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Performance Tips

1. **Use pagination** for large result sets
2. **Cache responses** when appropriate
3. **Use specific filters** to reduce response size
4. **Run scraper during off-peak hours** for better performance
5. **Adjust worker count** based on your system resources
