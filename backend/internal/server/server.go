package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"coderelay/backend/internal/storage"
)

// Config controls how the HTTP server listens for requests.
type Config struct {
	Addr   string
	DBPath string
}

// healthResponse is returned from the health check endpoint.
type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Server manages the HTTP listeners for REST and WebSocket traffic.
type Server struct {
	cfg        Config
	httpServer *http.Server
	store      *storage.SQLite
	ownsStore  bool // whether server should close the store
}

// New will wire up routing, storage, and worker coordination.
func New(cfg Config) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "./coderelay.db"
	}

	mux := http.NewServeMux()
	srv := &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.Addr,
			Handler: corsMiddleware(mux),
		},
		ownsStore: true,
	}

	// Health check
	mux.HandleFunc("GET /healthz", srv.handleHealth)

	// API routes
	mux.HandleFunc("GET /api/problems", srv.handleListProblems)
	mux.HandleFunc("GET /api/problems/", srv.handleGetProblem)
	mux.HandleFunc("POST /api/submissions", srv.handleCreateSubmission)
	mux.HandleFunc("GET /api/submissions/", srv.handleGetSubmission)

	return srv
}

// Start boots the HTTP server and blocks until the context is canceled.
func (s *Server) Start(ctx context.Context) error {
	// Initialize database if not already set
	if s.store == nil {
		store, err := storage.Open(s.cfg.DBPath)
		if err != nil {
			return err
		}
		s.store = store
		s.ownsStore = true

		// Run migrations and seed
		if err := s.store.Migrate(); err != nil {
			return err
		}
		if err := s.store.Seed(); err != nil {
			// Seed errors are non-fatal (data might already exist)
		}
	}

	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if s.store != nil && s.ownsStore {
			s.store.Close()
		}

		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "health encoding failed", http.StatusInternalServerError)
	}
}

// corsMiddleware adds CORS headers for development.
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
