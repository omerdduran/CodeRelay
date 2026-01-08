package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// --- Request/Response Types ---

// SubmissionRequest is the request body for creating a submission.
type SubmissionRequest struct {
	UserID    int64  `json:"user_id"`
	ProblemID int64  `json:"problem_id"`
	Code      string `json:"code"`
	Language  string `json:"language"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// LeaderboardEntry represents a player's ranking.
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	UserID      int64  `json:"user_id"`
	Nickname    string `json:"nickname"`
	Solved      bool   `json:"solved"`
	SolveTimeMs *int   `json:"solve_time_ms,omitempty"`
}

// --- Handlers ---

// handleListProblems returns all problems.
func (s *Server) handleListProblems(w http.ResponseWriter, r *http.Request) {
	problems, err := s.store.ListProblems()
	if err != nil {
		s.jsonError(w, "failed to list problems", http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, problems, http.StatusOK)
}

// handleGetProblem returns a single problem with sample test cases.
func (s *Server) handleGetProblem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/problems/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.jsonError(w, "invalid problem id", http.StatusBadRequest)
		return
	}

	problem, err := s.store.GetProblemByID(id)
	if err != nil {
		s.jsonError(w, "problem not found", http.StatusNotFound)
		return
	}

	s.jsonResponse(w, problem, http.StatusOK)
}

// handleCreateSubmission creates a new submission.
func (s *Server) handleCreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.ProblemID == 0 || req.Code == "" {
		s.jsonError(w, "user_id, problem_id, and code are required", http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}

	submission, err := s.store.CreateSubmission(req.UserID, req.ProblemID, req.Code, req.Language)
	if err != nil {
		s.jsonError(w, "failed to create submission", http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, submission, http.StatusCreated)
}

// handleGetSubmission returns a single submission.
func (s *Server) handleGetSubmission(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/submissions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.jsonError(w, "invalid submission id", http.StatusBadRequest)
		return
	}

	submission, err := s.store.GetSubmissionByID(id)
	if err != nil {
		s.jsonError(w, "submission not found", http.StatusNotFound)
		return
	}

	s.jsonResponse(w, submission, http.StatusOK)
}

// handleLeaderboard returns rankings - either ELO-based or problem-specific.
func (s *Server) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboardType := r.URL.Query().Get("type")
	
	// If type=elo or no problem_id specified, return ELO leaderboard
	if leaderboardType == "elo" || r.URL.Query().Get("problem_id") == "" {
		s.handleEloLeaderboard(w, r)
		return
	}

	// Otherwise return problem-specific leaderboard
	s.handleProblemLeaderboard(w, r)
}

// handleEloLeaderboard returns global ELO rankings.
func (s *Server) handleEloLeaderboard(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	entries, err := s.store.GetLeaderboard(limit)
	if err != nil {
		s.jsonError(w, "failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, entries, http.StatusOK)
}

// handleProblemLeaderboard returns rankings for a specific problem.
func (s *Server) handleProblemLeaderboard(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.URL.Query().Get("problem_id")
	problemID, err := strconv.ParseInt(problemIDStr, 10, 64)
	if err != nil {
		s.jsonError(w, "invalid problem_id", http.StatusBadRequest)
		return
	}

	rows, err := s.store.DB.Query(`
		SELECT s.user_id, u.nickname, s.runtime_ms
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.problem_id = ? AND s.status = 'AC'
		ORDER BY s.runtime_ms ASC, s.created_at ASC
	`, problemID)
	if err != nil {
		s.jsonError(w, "failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	seenUsers := make(map[int64]bool)
	rank := 0

	for rows.Next() {
		var userID int64
		var nickname string
		var runtimeMs int

		if err := rows.Scan(&userID, &nickname, &runtimeMs); err != nil {
			continue
		}

		if seenUsers[userID] {
			continue
		}
		seenUsers[userID] = true
		rank++

		entries = append(entries, LeaderboardEntry{
			Rank:        rank,
			UserID:      userID,
			Nickname:    nickname,
			Solved:      true,
			SolveTimeMs: &runtimeMs,
		})
	}

	s.jsonResponse(w, entries, http.StatusOK)
}

// UserRequest is the request body for creating a user.
type UserRequest struct {
	Nickname string `json:"nickname"`
}

// handleCreateUser creates or returns an existing user.
func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" {
		s.jsonError(w, "nickname is required", http.StatusBadRequest)
		return
	}

	// Try to find existing user
	row := s.store.DB.QueryRow("SELECT id, nickname FROM users WHERE nickname = ?", req.Nickname)
	var user struct {
		ID       int64  `json:"id"`
		Nickname string `json:"nickname"`
	}

	err := row.Scan(&user.ID, &user.Nickname)
	if err == nil {
		// User exists
		s.jsonResponse(w, user, http.StatusOK)
		return
	}

	// Create new user
	result, err := s.store.DB.Exec("INSERT INTO users (nickname) VALUES (?)", req.Nickname)
	if err != nil {
		s.jsonError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	user.ID, _ = result.LastInsertId()
	user.Nickname = req.Nickname

	s.jsonResponse(w, user, http.StatusCreated)
}

// handleGetUserInfo returns user info with ELO rating and history.
func (s *Server) handleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	
	// Check if requesting ELO history
	if strings.HasSuffix(path, "/elo-history") {
		userIDStr := strings.TrimSuffix(path, "/elo-history")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			s.jsonError(w, "invalid user id", http.StatusBadRequest)
			return
		}
		
		limitStr := r.URL.Query().Get("limit")
		limit := 50
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}
		
		history, err := s.store.GetUserEloHistory(userID, limit)
		if err != nil {
			s.jsonError(w, "failed to fetch ELO history", http.StatusInternalServerError)
			return
		}
		
		s.jsonResponse(w, history, http.StatusOK)
		return
	}
	
	// Otherwise return basic user info
	userID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		s.jsonError(w, "invalid user id", http.StatusBadRequest)
		return
	}
	
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		s.jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	
	s.jsonResponse(w, user, http.StatusOK)
}

// --- Helper Methods ---

func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
