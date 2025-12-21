package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"coderelay/backend/internal/storage"
)

func setupTestServer(t *testing.T) (*Server, func()) {
	dbPath := "./test_server.db"

	store, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := store.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	if err := store.Seed(); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	srv := New(Config{Addr: ":8080", DBPath: dbPath})
	srv.store = store

	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return srv, cleanup
}

func TestHandleListProblems(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/problems", nil)
	w := httptest.NewRecorder()

	srv.handleListProblems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var problems []storage.Problem
	if err := json.NewDecoder(w.Body).Decode(&problems); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %d", len(problems))
	}
}

func TestHandleGetProblem(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/problems/1", nil)
	w := httptest.NewRecorder()

	srv.handleGetProblem(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var problem storage.Problem
	if err := json.NewDecoder(w.Body).Decode(&problem); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if problem.Title != "Two Sum" {
		t.Errorf("expected 'Two Sum', got '%s'", problem.Title)
	}

	if len(problem.SampleCases) != 2 {
		t.Errorf("expected 2 sample cases, got %d", len(problem.SampleCases))
	}
}

func TestHandleGetProblemNotFound(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/problems/999", nil)
	w := httptest.NewRecorder()

	srv.handleGetProblem(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandleCreateSubmission(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := SubmissionRequest{
		UserID:    1,
		ProblemID: 1,
		Code:      "print('hello')",
		Language:  "python",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/submissions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleCreateSubmission(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var submission storage.Submission
	if err := json.NewDecoder(w.Body).Decode(&submission); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if submission.Status != "queued" {
		t.Errorf("expected status 'queued', got '%s'", submission.Status)
	}
}

func TestHandleCreateSubmissionValidation(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	// Missing required fields
	body := SubmissionRequest{
		UserID: 1,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/submissions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleCreateSubmission(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGetSubmission(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	// First create a submission
	_, err := srv.store.CreateSubmission(1, 1, "print('test')", "python")
	if err != nil {
		t.Fatalf("failed to create submission: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/submissions/1", nil)
	w := httptest.NewRecorder()

	srv.handleGetSubmission(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var submission storage.Submission
	if err := json.NewDecoder(w.Body).Decode(&submission); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if submission.Code != "print('test')" {
		t.Errorf("expected code 'print('test')', got '%s'", submission.Code)
	}
}
