package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	TypeSubmissionUpdate MessageType = "submission_update"
	TypeLeaderboard      MessageType = "leaderboard"
	TypeTimer            MessageType = "timer"
	TypeError            MessageType = "error"
	TypeRaceEvent        MessageType = "race_event"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// SubmissionUpdate is sent when a submission status changes
type SubmissionUpdate struct {
	SubmissionID int64  `json:"submission_id"`
	UserID       int64  `json:"user_id"`
	ProblemID    int64  `json:"problem_id"`
	Status       string `json:"status"`
	RuntimeMs    *int   `json:"runtime_ms,omitempty"`
}

// LeaderboardEntry represents a user's position
type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	Solved    bool   `json:"solved"`
	SolveTime *int   `json:"solve_time_ms,omitempty"`
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast channel for messages
	broadcast chan []byte

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("ws: client connected, total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("ws: client disconnected, total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msgType MessageType, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.broadcast <- msgBytes
	return nil
}

// BroadcastSubmissionUpdate sends a submission update to all clients
func (h *Hub) BroadcastSubmissionUpdate(subID, userID, problemID int64, status string, runtimeMs *int) {
	update := SubmissionUpdate{
		SubmissionID: subID,
		UserID:       userID,
		ProblemID:    problemID,
		Status:       status,
		RuntimeMs:    runtimeMs,
	}
	if err := h.Broadcast(TypeSubmissionUpdate, update); err != nil {
		log.Printf("ws: failed to broadcast submission update: %v", err)
	}
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
