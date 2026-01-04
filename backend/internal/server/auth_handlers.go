package server

import (
	"encoding/json"
	"net/http"

	"coderelay/backend/internal/auth"
)

// --- Auth Request/Response Types ---

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// AuthResponse is the response for successful auth operations.
type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64  `json:"id"`
		Nickname string `json:"nickname"`
	} `json:"user"`
}

// --- Auth Handlers ---

// handleRegister creates a new user with password.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" || req.Password == "" {
		s.jsonError(w, "nickname and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		s.jsonError(w, "nickname must be 2-20 characters", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		s.jsonError(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := s.store.GetUserByNickname(req.Nickname)
	if err == nil {
		s.jsonError(w, "nickname already taken", http.StatusConflict)
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		s.jsonError(w, "failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	var user interface{}
	if req.Email != "" {
		user, err = s.store.CreateUserWithEmail(req.Nickname, req.Email, passwordHash)
	} else {
		user, err = s.store.CreateUser(req.Nickname, passwordHash)
	}

	if err != nil {
		s.jsonError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// Get user struct to access ID
	storedUser, _ := s.store.GetUserByNickname(req.Nickname)

	// Generate token
	token, err := auth.GenerateToken(storedUser.ID, storedUser.Nickname)
	if err != nil {
		s.jsonError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Avoid unused variable
	_ = user

	response := AuthResponse{Token: token}
	response.User.ID = storedUser.ID
	response.User.Nickname = storedUser.Nickname

	s.jsonResponse(w, response, http.StatusCreated)
}

// handleLogin authenticates a user and returns a token.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" || req.Password == "" {
		s.jsonError(w, "nickname and password are required", http.StatusBadRequest)
		return
	}

	// Find user
	user, err := s.store.GetUserByNickname(req.Nickname)
	if err != nil {
		s.jsonError(w, "invalid nickname or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		s.jsonError(w, "invalid nickname or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Nickname)
	if err != nil {
		s.jsonError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token}
	response.User.ID = user.ID
	response.User.Nickname = user.Nickname

	s.jsonResponse(w, response, http.StatusOK)
}

// handleGetMe returns the current authenticated user.
func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		s.jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.store.GetUserByID(claims.UserID)
	if err != nil {
		s.jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	s.jsonResponse(w, map[string]interface{}{
		"id":       user.ID,
		"nickname": user.Nickname,
	}, http.StatusOK)
}
