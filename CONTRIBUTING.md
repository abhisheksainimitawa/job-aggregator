# Contributing to Job Aggregator

Thank you for your interest in contributing! 🎉

## How to Contribute

### Reporting Bugs
1. Check if the bug has already been reported in Issues
2. Create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version)

### Suggesting Features
1. Open an issue with the "enhancement" label
2. Describe the feature and its use case
3. Explain how it would benefit users

### Code Contributions

#### Setup Development Environment
```bash
# Fork and clone the repo
git clone https://github.com/abhisheksainimitawa/job-aggregator.git
cd job-aggregator

# Install dependencies
go mod download

# Start dependencies
docker-compose up -d postgres

# Run tests
go test ./...
```

#### Making Changes
1. Create a new branch: `git checkout -b feature/your-feature-name`
2. Make your changes
3. Write/update tests
4. Run tests: `go test ./...`
5. Format code: `go fmt ./...`
6. Commit with clear messages
7. Push and create a Pull Request

#### Code Style
- Follow Go conventions (`gofmt`, `golint`)
- Write clear comments for exported functions
- Keep functions focused and small
- Use meaningful variable names
- Add tests for new functionality

#### Pull Request Guidelines
- Link related issues
- Describe what changed and why
- Include test coverage
- Update documentation if needed
- Ensure CI passes

## Development Guidelines

### Adding a New Job Source
1. Implement the `JobSource` interface in `internal/scraper/sources.go`
2. Add tests in `internal/scraper/sources_test.go`
3. Register the source in `cmd/api/main.go` and `cmd/scraper/main.go`
4. Update documentation

### API Changes
- Maintain backward compatibility
- Version new endpoints appropriately
- Update EXAMPLES.md with usage

### Database Changes
- Create migration scripts
- Update schema documentation
- Test with existing data

## Questions?

Feel free to open an issue for any questions or clarifications!

---

Happy coding! 🚀
