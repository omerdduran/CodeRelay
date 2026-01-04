package server

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"coderelay/backend/internal/ws"
)

// RaceRequest is the request body for creating a race.
type RaceRequest struct {
	ProblemID int64 `json:"problem_id"`
	UserID    int64 `json:"user_id"`
}

// JoinRaceRequest is the request body for joining a race.
type JoinRaceRequest struct {
	UserID int64 `json:"user_id"`
}

// handleCreateRace creates a new race room.
func (s *Server) handleCreateRace(w http.ResponseWriter, r *http.Request) {
	var req RaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProblemID == 0 || req.UserID == 0 {
		s.jsonError(w, "problem_id and user_id are required", http.StatusBadRequest)
		return
	}

	// Generate room code
	code := generateRoomCode()

	race, err := s.store.CreateRace(code, req.ProblemID, req.UserID)
	if err != nil {
		s.jsonError(w, "failed to create race", http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, race, http.StatusCreated)
}

// handleGetRace returns a race by room code.
func (s *Server) handleGetRace(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/races/")
	code = strings.Split(code, "/")[0]

	if code == "" {
		s.jsonError(w, "room code required", http.StatusBadRequest)
		return
	}

	race, err := s.store.GetRaceByCode(code)
	if err != nil {
		s.jsonError(w, "race not found", http.StatusNotFound)
		return
	}

	s.jsonResponse(w, race, http.StatusOK)
}

// handleJoinRace adds a player to a race.
func (s *Server) handleJoinRace(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/races/")
	code = strings.TrimSuffix(code, "/join")

	var req JoinRaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	race, err := s.store.GetRaceByCode(code)
	if err != nil {
		s.jsonError(w, "race not found", http.StatusNotFound)
		return
	}

	if race.Status != "waiting" {
		s.jsonError(w, "race already started", http.StatusBadRequest)
		return
	}

	if err := s.store.JoinRace(race.ID, req.UserID); err != nil {
		s.jsonError(w, "failed to join race", http.StatusInternalServerError)
		return
	}

	// Get updated race
	race, _ = s.store.GetRaceByCode(code)

	// Broadcast player joined
	if s.hub != nil {
		s.hub.Broadcast(ws.TypeRaceEvent, map[string]interface{}{
			"event":     "player_joined",
			"room_code": code,
			"user_id":   req.UserID,
		})
	}

	s.jsonResponse(w, race, http.StatusOK)
}

// handleStartRace starts the race.
func (s *Server) handleStartRace(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/races/")
	code = strings.TrimSuffix(code, "/start")

	var req JoinRaceRequest
	json.NewDecoder(r.Body).Decode(&req)

	race, err := s.store.GetRaceByCode(code)
	if err != nil {
		s.jsonError(w, "race not found", http.StatusNotFound)
		return
	}

	if race.HostUserID != req.UserID {
		s.jsonError(w, "only host can start race", http.StatusForbidden)
		return
	}

	if race.Status != "waiting" {
		s.jsonError(w, "race already started", http.StatusBadRequest)
		return
	}

	// Start countdown
	if s.hub != nil {
		s.hub.Broadcast(ws.TypeRaceEvent, map[string]interface{}{
			"event":     "countdown",
			"room_code": code,
			"seconds":   3,
		})
	}

	// Wait for countdown then start
	go func() {
		time.Sleep(3 * time.Second)

		startTime := time.Now()
		s.store.UpdateRaceStatus(race.ID, "racing", &startTime)

		// Update all participants to racing
		for _, p := range race.Players {
			s.store.UpdateRaceParticipant(race.ID, p.UserID, "racing", nil, nil)
		}

		if s.hub != nil {
			s.hub.Broadcast(ws.TypeRaceEvent, map[string]interface{}{
				"event":      "race_started",
				"room_code":  code,
				"start_time": startTime,
			})
		}
	}()

	s.jsonResponse(w, map[string]string{"status": "countdown"}, http.StatusOK)
}

// generateRoomCode creates a 6-character room code.
func generateRoomCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

// handleWatchRace adds a spectator to a race (allowed anytime).
func (s *Server) handleWatchRace(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/races/")
	code = strings.TrimSuffix(code, "/watch")

	var req JoinRaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	race, err := s.store.GetRaceByCode(code)
	if err != nil {
		s.jsonError(w, "race not found", http.StatusNotFound)
		return
	}

	// Spectators can join anytime (waiting, racing, or finished)
	if err := s.store.JoinRaceAsSpectator(race.ID, req.UserID); err != nil {
		s.jsonError(w, "failed to join as spectator", http.StatusInternalServerError)
		return
	}

	// Get updated race
	race, _ = s.store.GetRaceByCode(code)

	// Broadcast spectator joined
	if s.hub != nil {
		s.hub.Broadcast(ws.TypeRaceEvent, map[string]interface{}{
			"event":     "spectator_joined",
			"room_code": code,
			"user_id":   req.UserID,
		})
	}

	s.jsonResponse(w, race, http.StatusOK)
}

// handleRaceAction routes POST requests to join, start, or watch handlers.
func (s *Server) handleRaceAction(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/join") {
		s.handleJoinRace(w, r)
	} else if strings.HasSuffix(path, "/start") {
		s.handleStartRace(w, r)
	} else if strings.HasSuffix(path, "/watch") {
		s.handleWatchRace(w, r)
	} else {
		s.jsonError(w, "invalid race action", http.StatusBadRequest)
	}
}
