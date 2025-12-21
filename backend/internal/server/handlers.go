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
