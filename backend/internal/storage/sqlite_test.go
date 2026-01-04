package storage

import (
	"os"
	"testing"
)

func TestSQLiteOpen(t *testing.T) {
	dbPath := "./test_coderelay.db"
	defer os.Remove(dbPath)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if store.DB == nil {
		t.Error("expected DB to be non-nil")
	}
}

func TestSQLiteMigrateAndSeed(t *testing.T) {
	dbPath := "./test_migrate.db"
	defer os.Remove(dbPath)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if err := store.Seed(); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	// Verify seeded data
	problems, err := store.ListProblems()
	if err != nil {
		t.Fatalf("failed to list problems: %v", err)
	}
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %d", len(problems))
	}
	if problems[0].Title != "Two Sum" {
		t.Errorf("expected 'Two Sum', got '%s'", problems[0].Title)
	}
}

func TestUserOperations(t *testing.T) {
	dbPath := "./test_users.db"
	defer os.Remove(dbPath)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Create user
	user, err := store.CreateUser("testplayer", "hashedpassword123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if user.Nickname != "testplayer" {
		t.Errorf("expected nickname 'testplayer', got '%s'", user.Nickname)
	}

	// Get user by ID
	fetched, err := store.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}
	if fetched.Nickname != "testplayer" {
		t.Errorf("expected nickname 'testplayer', got '%s'", fetched.Nickname)
	}

	// Get user by nickname
	fetched, err = store.GetUserByNickname("testplayer")
	if err != nil {
		t.Fatalf("failed to get user by nickname: %v", err)
	}
	if fetched.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, fetched.ID)
	}
}

func TestProblemOperations(t *testing.T) {
	dbPath := "./test_problems.db"
	defer os.Remove(dbPath)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	if err := store.Seed(); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	// Get problem by ID
	problem, err := store.GetProblemByID(1)
	if err != nil {
		t.Fatalf("failed to get problem: %v", err)
	}
	if problem.Title != "Two Sum" {
		t.Errorf("expected 'Two Sum', got '%s'", problem.Title)
	}

	// Should have sample test cases
	if len(problem.SampleCases) != 2 {
		t.Errorf("expected 2 sample cases, got %d", len(problem.SampleCases))
	}

	// Get all test cases (including hidden)
	cases, err := store.GetTestCasesByProblemID(1)
	if err != nil {
		t.Fatalf("failed to get test cases: %v", err)
	}
	if len(cases) != 5 {
		t.Errorf("expected 5 test cases, got %d", len(cases))
	}
}

func TestSubmissionOperations(t *testing.T) {
	dbPath := "./test_submissions.db"
	defer os.Remove(dbPath)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	if err := store.Seed(); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	// Create submission
	submission, err := store.CreateSubmission(1, 1, "print('hello')", "python")
	if err != nil {
		t.Fatalf("failed to create submission: %v", err)
	}
	if submission.Status != "queued" {
		t.Errorf("expected status 'queued', got '%s'", submission.Status)
	}

	// Update submission status
	runtime := 150
	if err := store.UpdateSubmissionStatus(submission.ID, "AC", &runtime); err != nil {
		t.Fatalf("failed to update submission: %v", err)
	}

	// Verify update
	updated, err := store.GetSubmissionByID(submission.ID)
	if err != nil {
		t.Fatalf("failed to get submission: %v", err)
	}
	if updated.Status != "AC" {
		t.Errorf("expected status 'AC', got '%s'", updated.Status)
	}
	if updated.RuntimeMs == nil || *updated.RuntimeMs != 150 {
		t.Errorf("expected runtime 150, got %v", updated.RuntimeMs)
	}
}
