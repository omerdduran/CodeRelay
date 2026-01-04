package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"coderelay/backend/internal/storage"
	"coderelay/backend/internal/ws"
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
	hub        *ws.Hub
	ownsStore  bool
}

// New will wire up routing, storage, and worker coordination.
func New(cfg Config, hub *ws.Hub) *Server {
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
		hub:       hub,
		ownsStore: true,
	}

	// Health check
	mux.HandleFunc("GET /healthz", srv.handleHealth)

	// API routes
	mux.HandleFunc("GET /api/problems", srv.handleListProblems)
	mux.HandleFunc("GET /api/problems/", srv.handleGetProblem)
	mux.HandleFunc("POST /api/submissions", srv.handleCreateSubmission)
	mux.HandleFunc("GET /api/submissions/", srv.handleGetSubmission)
	mux.HandleFunc("GET /api/leaderboard", srv.handleLeaderboard)
	mux.HandleFunc("POST /api/users", srv.handleCreateUser)

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", srv.handleRegister)
	mux.HandleFunc("POST /api/auth/login", srv.handleLogin)
	mux.HandleFunc("GET /api/auth/me", srv.handleGetMe)

	// Race routes
	mux.HandleFunc("POST /api/races", srv.handleCreateRace)
	mux.HandleFunc("GET /api/races/", srv.handleGetRace)
	mux.HandleFunc("POST /api/races/", srv.handleRaceAction)

	// WebSocket
	mux.HandleFunc("GET /ws", srv.handleWebSocket)

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

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(s.hub, w, r)
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
