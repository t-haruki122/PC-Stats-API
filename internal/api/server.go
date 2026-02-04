package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"pc-stats-api/internal/storage"
)

// Server represents the HTTP API server
type Server struct {
	storage  *storage.RingBuffer
	interval int
	mux      *http.ServeMux
}

// NewServer creates a new API server
func NewServer(storage *storage.RingBuffer, intervalSec int) *Server {
	s := &Server{
		storage:  storage,
		interval: intervalSec,
		mux:      http.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// API endpoints
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/metrics/latest", s.handleLatest)
	s.mux.HandleFunc("/metrics/history", s.handleHistory)

	// Static files for Web UI
	fs := http.FileServer(http.Dir("web"))
	s.mux.Handle("/ui/", http.StripPrefix("/ui/", fs))
	s.mux.HandleFunc("/", s.handleRoot)
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	addr := ":" + port
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, s.enableCORS(s.mux))
}

// enableCORS adds CORS headers
func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleRoot redirects to Web UI
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/ui/", http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// handleLatest returns the most recent metrics
func (s *Server) handleLatest(w http.ResponseWriter, r *http.Request) {
	latest := s.storage.GetLatest()
	if latest == nil {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latest)
}

// handleHistory returns historical metrics
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	// Parse seconds parameter (default: 300)
	secondsStr := r.URL.Query().Get("seconds")
	seconds := 300
	if secondsStr != "" {
		if parsed, err := strconv.Atoi(secondsStr); err == nil {
			seconds = parsed
		}
	}

	history := s.storage.GetHistory(seconds)

	response := map[string]interface{}{
		"interval_sec": s.interval,
		"samples":      history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
