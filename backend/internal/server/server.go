package server

import "context"

// Server represents the future HTTP/WebSocket server stack.
type Server struct{}

// New will wire up routing, storage, and worker coordination.
func New() *Server {
	return &Server{}
}

// Start blocks until the provided context is canceled; real listeners land later.
func (s *Server) Start(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
