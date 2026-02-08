package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/abhisheksainimitawa/job-aggregator/internal/models"
	"github.com/abhisheksainimitawa/job-aggregator/internal/service"
	"github.com/abhisheksainimitawa/job-aggregator/pkg/logger"
)

// Handler holds all HTTP handlers
type Handler struct {
	jobService *service.JobService
}

// NewHandler creates a new HTTP handler
func NewHandler(jobService *service.JobService) *Handler {
	return &Handler{
		jobService: jobService,
	}
}

// SetupRoutes sets up all API routes
func (h *Handler) SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Job routes
	api.HandleFunc("/jobs", h.ListJobs).Methods("GET")
	api.HandleFunc("/jobs/{id:[0-9]+}", h.GetJob).Methods("GET")
	api.HandleFunc("/jobs/search", h.SearchJobs).Methods("GET")
	api.HandleFunc("/jobs/stats", h.GetStats).Methods("GET")

	// Scraper routes
	api.HandleFunc("/scraper/run", h.RunScraper).Methods("POST")
	api.HandleFunc("/scraper/status", h.GetScraperStatus).Methods("GET")

	// Health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Apply middleware
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	return r
}

// ListJobs lists all jobs with pagination
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20
	}

	query := &models.JobSearchQuery{
		Page:  page,
		Limit: limit,
	}

	jobs, err := h.jobService.SearchJobs(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch jobs")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":  jobs,
		"page":  page,
		"limit": limit,
		"total": len(jobs),
	})
}

// GetJob retrieves a single job by ID
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.jobService.GetJob(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	respondJSON(w, http.StatusOK, job)
}

// SearchJobs searches for jobs based on query parameters
func (h *Handler) SearchJobs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := &models.JobSearchQuery{
		Keywords: q.Get("q"),
		Location: q.Get("location"),
		JobType:  q.Get("type"),
		Source:   q.Get("source"),
		Page:     0,
		Limit:    20,
	}

	if page := q.Get("page"); page != "" {
		query.Page, _ = strconv.Atoi(page)
	}

	if limit := q.Get("limit"); limit != "" {
		query.Limit, _ = strconv.Atoi(limit)
	}

	if remote := q.Get("remote"); remote == "true" {
		t := true
		query.Remote = &t
	} else if remote == "false" {
		f := false
		query.Remote = &f
	}

	jobs, err := h.jobService.SearchJobs(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Search failed")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":   jobs,
		"query":  query.Keywords,
		"total":  len(jobs),
		"page":   query.Page,
		"limit":  query.Limit,
	})
}

// GetStats retrieves job statistics
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.jobService.GetStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// RunScraper triggers the scraper
func (h *Handler) RunScraper(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Query == "" {
		req.Query = "golang developer"
	}

	// Run scraper with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	count, err := h.jobService.RunScraper(ctx, req.Query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Scraper failed: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Scraper completed successfully",
		"jobs_scraped": count,
		"query":        req.Query,
	})
}

// GetScraperStatus returns the current scraper status
func (h *Handler) GetScraperStatus(w http.ResponseWriter, r *http.Request) {
	stats := h.jobService.GetScraperStats()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"jobs_scraped": stats.JobsScraped,
		"errors":       stats.Errors,
		"start_time":   stats.StartTime,
		"end_time":     stats.EndTime,
	})
}

// HealthCheck returns API health status
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "job-aggregator",
		"time":    time.Now(),
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Middleware

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		logger.Debug("Request completed in %v", time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
