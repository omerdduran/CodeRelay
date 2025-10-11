package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Config controls how the HTTP server listens for requests.
type Config struct {
	Addr string
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
}

// New will wire up routing, storage, and worker coordination.
func New(cfg Config) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}

	mux := http.NewServeMux()
	srv := &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
	}

	mux.HandleFunc("GET /healthz", srv.handleHealth)

	return srv
}

// Start boots the HTTP server and blocks until the context is canceled.
func (s *Server) Start(ctx context.Context) error {
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
