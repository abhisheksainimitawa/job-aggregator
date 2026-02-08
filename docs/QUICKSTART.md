# 🚀 Quick Start Guide - 5 Minutes to Running

## Prerequisites Check

✅ **Go installed?** Run: `go version` (need 1.21+)  
✅ **Docker installed?** Run: `docker --version`

If not installed:
- Go: https://golang.org/dl/
- Docker: https://www.docker.com/products/docker-desktop/

---

## ⚡ Fast Track (3 Steps)

### Step 1: Install Go Dependencies (30 seconds)
```bash
# Navigate to project directory
cd job-aggregator  # or wherever you cloned/downloaded the project
go mod download
```

### Step 2: Start Database & Redis (1 minute)
```bash
docker-compose up -d postgres redis
```

Wait 10 seconds for services to start, then verify:
```bash
docker-compose ps
```

You should see postgres and redis running.

### Step 3: Run the Application (10 seconds)

**Option A: Start API Server**
```bash
go run cmd/api/main.go
```

**Option B: Run Scraper**
```bash
go run cmd/scraper/main.go
```

The API will be available at `http://localhost:8080`

---

## 📋 Environment Configuration (Optional)

The `.env` file is already created with these defaults:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=job_aggregator
DB_SSLMODE=disable

# Redis Configuration  
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Scraper Configuration
SCRAPER_WORKERS=10
SCRAPER_RATE_LIMIT=100
SCRAPER_TIMEOUT=30
```

You can modify these values if needed.

---

## ✅ Test It Works

### Test API Server
Open a new terminal and run:
```bash
# Health check
curl http://localhost:8080/health

# Trigger scraping
curl -X POST http://localhost:8080/api/v1/scraper/run -H "Content-Type: application/json" -d "{\"query\": \"golang developer\"}"

# View jobs
curl http://localhost:8080/api/v1/jobs

# Search for specific jobs
curl "http://localhost:8080/api/v1/jobs/search?q=backend&location=remote"

# Get statistics
curl http://localhost:8080/api/v1/jobs/stats
```

### Test Scraper CLI
```bash
go run cmd/scraper/main.go -query "backend engineer"
```

You should see:
```
[INFO] Starting Job Scraper CLI...
[INFO] Registered scraper source: Indeed
[INFO] Registered scraper source: LinkedIn
[INFO] Registered scraper source: Glassdoor
[INFO] Starting scraper engine with 10 workers...
[INFO] Scraping completed: 25 jobs, 0 errors, duration: 4.2s
```

---

## 🎯 Your First Task

Try this challenge:
1. Start the API server
2. Trigger a scrape via API
3. Search for "remote" jobs
4. Get statistics

Solution:
```bash
# Terminal 1: Start server
go run cmd/api/main.go

# Terminal 2: Use the API
curl -X POST http://localhost:8080/api/v1/scraper/run -H "Content-Type: application/json" -d "{\"query\":\"golang\"}"
curl "http://localhost:8080/api/v1/jobs/search?remote=true"
curl http://localhost:8080/api/v1/jobs/stats
```

---

## 🐛 Troubleshooting

### "go: command not found"
→ Install Go from https://golang.org/dl/  
→ Restart your terminal after installation

### "Cannot connect to database"
→ Make sure Docker is running  
→ Run: `docker-compose up -d postgres`  
→ Wait 10-15 seconds for PostgreSQL to start

### "Redis connection failed"
→ Ensure Redis is running  
→ Run: `docker-compose up -d redis`

### "Port 8080 already in use"
→ Stop other programs using port 8080  
→ Or change SERVER_PORT in .env file

### "Docker not running"
→ Start Docker Desktop  
→ Wait for it to fully start (whale icon in taskbar)

---

## 📚 What to Do Next

1. **Explore the API** - Check [EXAMPLES.md](EXAMPLES.md)
2. **Read the code** - Start with [cmd/api/main.go](cmd/api/main.go)
3. **Run tests** - `go test ./...`
4. **Customize** - Add your own job sources
5. **Deploy** - Build Docker image: `docker-compose up --build`

---

## 🎓 Learning Path

**Day 1:** Understand the architecture ([ARCHITECTURE.md](ARCHITECTURE.md))  
**Day 2:** Study concurrency implementation ([internal/scraper/engine.go](internal/scraper/engine.go))  
**Day 3:** Explore database layer ([internal/repository/](internal/repository/))  
**Day 4:** Review API design ([internal/api/handler.go](internal/api/handler.go))  
**Day 5:** Write tests and add new features

---

## 🔑 Key Files to Know

| File | What It Does |
|------|--------------|
| `cmd/api/main.go` | Starts the HTTP server |
| `cmd/scraper/main.go` | CLI tool to scrape jobs |
| `internal/scraper/engine.go` | Concurrent scraping with goroutines |
| `internal/api/handler.go` | REST API endpoints |
| `internal/repository/job_repository.go` | Database operations |
| `pkg/ratelimit/ratelimit.go` | Rate limiting algorithm |

---

## 💡 Quick Commands Reference

```bash
# Start everything with Docker
docker-compose up --build

# Run API server
go run cmd/api/main.go

# Run scraper
go run cmd/scraper/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Build binaries (Linux/Mac)
go build -o bin/api cmd/api/main.go
go build -o bin/scraper cmd/scraper/main.go

# Build binaries (Windows)
go build -o bin/api.exe cmd/api/main.go
go build -o bin/scraper.exe cmd/scraper/main.go

# Stop Docker services
docker-compose down

# View Docker logs
docker-compose logs -f
```

---

## 🐳 Full Docker Deployment

Build and run everything with Docker:

```bash
# Build and start all services
docker-compose up --build

# Or in background
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop everything
docker-compose down
```

---

## 🔧 Development Workflow

1. **Make changes to code**
2. **Run tests:**
   ```bash
   go test ./...
   ```
3. **Format code:**
   ```bash
   go fmt ./...
   ```
4. **Build:**
   ```bash
   go build -o bin/api cmd/api/main.go
   go build -o bin/scraper cmd/scraper/main.go
   ```
5. **Run:**
   ```bash
   ./bin/api
   # or
   ./bin/scraper -query "your query"
   ```

---

## 📂 Project Structure Overview

```
job-aggregator/
├── cmd/
│   ├── api/           # API server (main entry point)
│   └── scraper/       # Scraper CLI tool
├── internal/
│   ├── api/           # HTTP handlers & routing
│   ├── service/       # Business logic
│   ├── repository/    # Database access
│   ├── scraper/       # Scraping engine
│   ├── models/        # Data models
│   └── config/        # Configuration
├── pkg/
│   ├── logger/        # Logging utility
│   └── ratelimit/     # Rate limiting
├── docker-compose.yml # Docker orchestration
├── Dockerfile         # Container image
├── go.mod            # Go dependencies
└── README.md         # Documentation
```

---

## ⚙️ Key Features Demonstrated

✅ **Goroutines & Channels** - Concurrent scraping with worker pools  
✅ **Context** - Timeout and cancellation handling  
✅ **Interfaces** - Clean abstraction for scrapers  
✅ **Database** - PostgreSQL integration with connection pooling  
✅ **REST API** - Production-ready HTTP server  
✅ **Rate Limiting** - Token bucket algorithm  
✅ **Error Handling** - Comprehensive error management  
✅ **Testing** - Unit tests for core components  
✅ **Docker** - Containerization and orchestration  
✅ **Logging** - Structured logging  
✅ **Configuration** - Environment-based config  

---

## 🚀 Next Steps / Enhancements

- [ ] Add Redis caching layer
- [ ] Implement real scraping (Colly/HTTP client)
- [ ] Add authentication (JWT)
- [ ] Build React frontend dashboard
- [ ] Add Prometheus metrics
- [ ] Implement job alerts via email
- [ ] Add Elasticsearch for advanced search
- [ ] Create Kubernetes manifests
- [ ] Add machine learning recommendations

---

## 🎉 You're Ready!

You now have a working Go project that demonstrates:
- ✅ Concurrent programming
- ✅ REST APIs
- ✅ Database integration
- ✅ Production patterns

**Next:** Share it on GitHub and start building! 🚀

---

## 📞 Support

For issues or questions:
- Check [EXAMPLES.md](EXAMPLES.md) for API usage examples
- Review the code documentation
- Test individual components with `go test ./...`

---

**Happy Coding! 🎉**
