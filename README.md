# 🚀 Job Aggregator

A high-performance job aggregation platform built with Go that scrapes multiple job boards concurrently and provides a unified search API.

## ✨ Features

- **Concurrent scraping** with configurable worker pools and goroutines
- **Intelligent deduplication** using SHA-256 hashing
- **RESTful API** with search, filtering, and statistics endpoints
- **Rate limiting** with token bucket algorithm
- **Production-ready** with error handling, logging, and graceful shutdown
- **Docker support** for easy deployment

## 🛠️ Tech Stack

- **Go 1.21+** - Core language
- **Gorilla Mux** - HTTP routing
- **PostgreSQL** - Database with connection pooling
- **Redis** - Caching layer
- **Docker** - Containerization

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Installation

```bash
# Clone the repository
git clone https://github.com/abhisheksainimitawa/job-aggregator.git
cd job-aggregator

# Install dependencies
go mod download

# Start database and Redis
docker-compose up -d postgres redis

# Run the API server
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Run the Scraper

```bash
go run cmd/scraper/main.go -query "golang developer"
```

## 📡 API Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Trigger scraping
curl -X POST http://localhost:8080/api/v1/scraper/run \
  -H "Content-Type: application/json" \
  -d '{"query": "golang developer"}'

# Search jobs
curl "http://localhost:8080/api/v1/jobs/search?q=backend&location=remote"

# Get statistics
curl http://localhost:8080/api/v1/jobs/stats
```

## 📂 Project Structure

```
job-aggregator/
├── cmd/
│   ├── api/           # API server
│   └── scraper/       # CLI scraper tool
├── internal/
│   ├── api/           # HTTP handlers
│   ├── service/       # Business logic
│   ├── repository/    # Database layer
│   ├── scraper/       # Scraping engine
│   ├── models/        # Data models
│   └── config/        # Configuration
├── pkg/
│   ├── logger/        # Logging utility
│   └── ratelimit/     # Rate limiter
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## 🔑 Key Concepts Demonstrated

- **Goroutines & Channels** - Concurrent worker pools
- **Context** - Timeout and cancellation handling
- **Interfaces** - Clean abstraction for scrapers
- **Clean Architecture** - Separation of concerns
- **Error Handling** - Graceful degradation
- **Testing** - Unit tests for core components

## 🐳 Docker Deployment

```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/scraper/...
```

## 📚 Documentation

- [QUICKSTART.md](docs/QUICKSTART.md) - Detailed setup guide
- [API_EXAMPLES.md](docs/API_EXAMPLES.md) - API usage examples
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - System design and patterns
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

## 🔧 Configuration

Edit `.env` file to configure:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=job_aggregator

# Server
SERVER_PORT=8080

# Scraper
SCRAPER_WORKERS=10
SCRAPER_RATE_LIMIT=100
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

Abhishek Saini - [GitHub](https://github.com/abhisheksainimitawa)

---

⭐ Star this repo if you find it useful!
