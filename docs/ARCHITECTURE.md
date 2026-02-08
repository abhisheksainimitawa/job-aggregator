# Architecture & Design

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                          │
│  (Browser, CLI, Postman, curl, Third-party Apps)            │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ HTTP/REST
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                      API Layer                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  HTTP Handlers (Gorilla Mux)                         │   │
│  │  - Jobs Handler                                      │   │
│  │  - Scraper Handler                                   │   │
│  │  - Stats Handler                                     │   │
│  └────────────────────┬─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                        │
                        │ Function Calls
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                   Service Layer                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Business Logic                                      │   │
│  │  - JobService (orchestrates scraping & storage)     │   │
│  │  - Validation & Deduplication                       │   │
│  └────────────────────┬─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                        │
        ┌───────────────┴──────────────┐
        │                              │
        ▼                              ▼
┌──────────────────┐          ┌──────────────────┐
│  Repository      │          │  Scraper Engine  │
│  Layer           │          │                  │
│                  │          │  Worker Pool     │
│  - CRUD Ops      │          │  (Goroutines)    │
│  - Queries       │          │                  │
│  - Batch Insert  │          │  ┌────────────┐  │
│                  │          │  │ Worker 1   │  │
└────────┬─────────┘          │  ├────────────┤  │
         │                    │  │ Worker 2   │  │
         │                    │  ├────────────┤  │
         ▼                    │  │ Worker 3   │  │
┌──────────────────┐          │  ├────────────┤  │
│   PostgreSQL     │          │  │ ...        │  │
│   Database       │          │  └────────────┘  │
│                  │          │                  │
│  - jobs table    │          │  Rate Limiter    │
│  - indexes       │          └────────┬─────────┘
└──────────────────┘                   │
                                       │
                     ┌─────────────────┼─────────────────┐
                     │                 │                 │
                     ▼                 ▼                 ▼
              ┌───────────┐     ┌───────────┐    ┌───────────┐
              │  Indeed   │     │ LinkedIn  │    │Glassdoor  │
              │  Scraper  │     │  Scraper  │    │ Scraper   │
              └───────────┘     └───────────┘    └───────────┘
```

## Concurrency Model

```
Main Goroutine
     │
     ├─── HTTP Server (goroutine per request)
     │
     └─── Scraper Engine
            │
            ├─── Worker 1 ────► Source Channel ────► Scrape Indeed
            │                                        ↓
            ├─── Worker 2 ────► Source Channel ────► Scrape LinkedIn
            │                                        ↓
            ├─── Worker 3 ────► Source Channel ────► Scrape Glassdoor
            │                                        ↓
            ├─── Worker 4 ────► Source Channel       │
            │                                        │
            └─── ...                                 │
                                                     ↓
                                           Jobs Channel
                                                     ↓
                                           Collector Goroutine
                                                     ↓
                                           Batch Insert to DB
```

## Data Flow

```
1. API Request
   ↓
2. Handler validates & routes
   ↓
3. Service coordinates business logic
   ↓
4. Scraper Engine dispatches work to pool
   ↓
5. Workers scrape sources concurrently
   ↓
6. Jobs sent through channel
   ↓
7. Collector aggregates results
   ↓
8. Repository deduplicates & stores
   ↓
9. Response sent to client
```

## Database Schema

```sql
┌─────────────────────────────────────────┐
│              jobs                       │
├─────────────────────────────────────────┤
│ id              BIGSERIAL PRIMARY KEY   │
│ title           VARCHAR(255) NOT NULL   │
│ company         VARCHAR(255) NOT NULL   │
│ location        VARCHAR(255) NOT NULL   │
│ salary          VARCHAR(100)            │
│ description     TEXT NOT NULL           │
│ url             TEXT NOT NULL           │
│ source          VARCHAR(50) NOT NULL    │
│ remote_ok       BOOLEAN                 │
│ job_type        VARCHAR(50) NOT NULL    │
│ posted_at       TIMESTAMP NOT NULL      │
│ scraped_at      TIMESTAMP NOT NULL      │
│ hash            VARCHAR(64) UNIQUE ◄──── Deduplication
│ created_at      TIMESTAMP NOT NULL      │
│ updated_at      TIMESTAMP NOT NULL      │
└─────────────────────────────────────────┘

Indexes:
- idx_jobs_title (title)
- idx_jobs_company (company)
- idx_jobs_location (location)
- idx_jobs_source (source)
- idx_jobs_posted_at (posted_at DESC)
- idx_jobs_remote_ok (remote_ok)
- idx_jobs_job_type (job_type)
- idx_jobs_hash (hash) - UNIQUE
```

## Key Design Patterns

### 1. Repository Pattern
Abstracts data access logic from business logic
```
Service → Repository Interface → Concrete Implementation → Database
```

### 2. Worker Pool Pattern
Limits concurrent operations and manages resources
```
Source Channel → Workers → Jobs Channel → Collector
```

### 3. Dependency Injection
Services receive dependencies through constructors
```
main.go creates all dependencies and injects them
```

### 4. Middleware Chain
Cross-cutting concerns handled uniformly
```
Request → Logging → CORS → Handler → Response
```

## Scalability Considerations

### Horizontal Scaling
- **Stateless API**: Can run multiple instances behind load balancer
- **Database Connection Pooling**: Efficient resource usage
- **Worker Pool**: Adjustable based on load

### Performance Optimizations
1. **Batch Inserts**: Reduce DB round trips
2. **Indexing**: Fast queries on common fields
3. **Rate Limiting**: Prevent overwhelming external sites
4. **Graceful Shutdown**: Clean resource cleanup

### Future Enhancements for Scale
- Redis caching for frequently accessed data
- Message queue (RabbitMQ/Kafka) for async processing
- Elasticsearch for advanced search
- CDN for static assets
- Kubernetes for orchestration

## Error Handling Strategy

```
1. Panic Recovery (middleware level)
   ↓
2. Structured Errors (domain level)
   ↓
3. Error Logging (observability)
   ↓
4. Graceful Degradation (continue on non-critical errors)
   ↓
5. User-Friendly Messages (API responses)
```

## Security Considerations

- Input validation on all API endpoints
- SQL injection prevention (parameterized queries)
- Rate limiting to prevent abuse
- CORS configuration for browser security
- Environment-based secrets (no hardcoded credentials)
- Future: JWT authentication, API keys

## Monitoring & Observability (Recommended)

```
Application
    ├── Structured Logs → Loki/ELK
    ├── Metrics → Prometheus
    └── Traces → Jaeger

Dashboards
    └── Grafana
```

## Testing Strategy

```
Unit Tests
├── Repository Layer (database mocking)
├── Service Layer (repository mocking)
├── Scraper Engine (source mocking)
└── Rate Limiter (time-based tests)

Integration Tests
└── Full API flow with test database

Performance Tests
└── Load testing with k6/Artillery
```

This architecture demonstrates production-level Go development practices!
