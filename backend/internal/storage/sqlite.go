package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

//go:embed seed.sql
var seedSQL string

// User represents a player in the system.
type User struct {
	ID           int64     `json:"id"`
	Nickname     string    `json:"nickname"`
	Email        *string   `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// Problem represents a coding challenge.
type Problem struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	TimeLimitMs   int        `json:"time_limit_ms"`
	MemoryLimitMb int        `json:"memory_limit_mb"`
	CreatedAt     time.Time  `json:"created_at"`
	SampleCases   []TestCase `json:"sample_cases,omitempty"`
}

// TestCase represents a test case for a problem.
type TestCase struct {
	ID             int64  `json:"id"`
	ProblemID      int64  `json:"problem_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       bool   `json:"is_sample"`
}

// Submission represents a code submission.
type Submission struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProblemID int64     `json:"problem_id"`
	Code      string    `json:"code"`
	Language  string    `json:"language"`
	Status    string    `json:"status"`
	RuntimeMs *int      `json:"runtime_ms,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// SQLite is a lightweight holder for the shared database handle.
type SQLite struct {
	DB *sql.DB
}

// Open creates or opens the SQLite database at the given path.
func Open(dbPath string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &SQLite{DB: db}, nil
}

// Migrate runs the schema migrations.
func (s *SQLite) Migrate() error {
	_, err := s.DB.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

// Seed populates the database with initial data.
func (s *SQLite) Seed() error {
	_, err := s.DB.Exec(seedSQL)
	if err != nil {
		return fmt.Errorf("run seed: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (s *SQLite) Close() error {
	return s.DB.Close()
}

// --- User Operations ---

// CreateUser creates a new user with the given nickname and password hash.
func (s *SQLite) CreateUser(nickname, passwordHash string) (*User, error) {
	result, err := s.DB.Exec("INSERT INTO users (nickname, password_hash) VALUES (?, ?)", nickname, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}

	return s.GetUserByID(id)
}

// CreateUserWithEmail creates a new user with nickname, email, and password hash.
func (s *SQLite) CreateUserWithEmail(nickname, email, passwordHash string) (*User, error) {
	result, err := s.DB.Exec("INSERT INTO users (nickname, email, password_hash) VALUES (?, ?, ?)", nickname, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}

	return s.GetUserByID(id)
}

// GetUserByID retrieves a user by ID.
func (s *SQLite) GetUserByID(id int64) (*User, error) {
	var u User
	var email sql.NullString
	err := s.DB.QueryRow(
		"SELECT id, nickname, email, password_hash, created_at FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Nickname, &email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if email.Valid {
		u.Email = &email.String
	}
	return &u, nil
}

// GetUserByNickname retrieves a user by nickname.
func (s *SQLite) GetUserByNickname(nickname string) (*User, error) {
	var u User
	var email sql.NullString
	err := s.DB.QueryRow(
		"SELECT id, nickname, email, password_hash, created_at FROM users WHERE nickname = ?", nickname,
	).Scan(&u.ID, &u.Nickname, &email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by nickname: %w", err)
	}
	if email.Valid {
		u.Email = &email.String
	}
	return &u, nil
}

// --- Problem Operations ---

// ListProblems returns all problems without test cases.
func (s *SQLite) ListProblems() ([]Problem, error) {
	rows, err := s.DB.Query(
		"SELECT id, title, description, time_limit_ms, memory_limit_mb, created_at FROM problems ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("list problems: %w", err)
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.TimeLimitMs, &p.MemoryLimitMb, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan problem: %w", err)
		}
		problems = append(problems, p)
	}
	return problems, nil
}

// GetProblemByID retrieves a problem by ID with sample test cases.
func (s *SQLite) GetProblemByID(id int64) (*Problem, error) {
	var p Problem
	err := s.DB.QueryRow(
		"SELECT id, title, description, time_limit_ms, memory_limit_mb, created_at FROM problems WHERE id = ?", id,
	).Scan(&p.ID, &p.Title, &p.Description, &p.TimeLimitMs, &p.MemoryLimitMb, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get problem by id: %w", err)
	}

	// Get sample test cases
	rows, err := s.DB.Query(
		"SELECT id, problem_id, input, expected_output, is_sample FROM test_cases WHERE problem_id = ? AND is_sample = 1", id,
	)
	if err != nil {
		return nil, fmt.Errorf("get sample test cases: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tc TestCase
		if err := rows.Scan(&tc.ID, &tc.ProblemID, &tc.Input, &tc.ExpectedOutput, &tc.IsSample); err != nil {
			return nil, fmt.Errorf("scan test case: %w", err)
		}
		p.SampleCases = append(p.SampleCases, tc)
	}

	return &p, nil
}

// GetTestCasesByProblemID returns all test cases for a problem (including hidden ones).
func (s *SQLite) GetTestCasesByProblemID(problemID int64) ([]TestCase, error) {
	rows, err := s.DB.Query(
		"SELECT id, problem_id, input, expected_output, is_sample FROM test_cases WHERE problem_id = ? ORDER BY id", problemID,
	)
	if err != nil {
		return nil, fmt.Errorf("get test cases: %w", err)
	}
	defer rows.Close()

	var cases []TestCase
	for rows.Next() {
		var tc TestCase
		if err := rows.Scan(&tc.ID, &tc.ProblemID, &tc.Input, &tc.ExpectedOutput, &tc.IsSample); err != nil {
			return nil, fmt.Errorf("scan test case: %w", err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// --- Submission Operations ---

// CreateSubmission creates a new submission.
func (s *SQLite) CreateSubmission(userID, problemID int64, code, language string) (*Submission, error) {
	result, err := s.DB.Exec(
		"INSERT INTO submissions (user_id, problem_id, code, language, status) VALUES (?, ?, ?, ?, 'queued')",
		userID, problemID, code, language,
	)
	if err != nil {
		return nil, fmt.Errorf("create submission: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}

	return s.GetSubmissionByID(id)
}

// GetSubmissionByID retrieves a submission by ID.
func (s *SQLite) GetSubmissionByID(id int64) (*Submission, error) {
	var sub Submission
	var runtimeMs sql.NullInt64
	err := s.DB.QueryRow(
		"SELECT id, user_id, problem_id, code, language, status, runtime_ms, created_at FROM submissions WHERE id = ?", id,
	).Scan(&sub.ID, &sub.UserID, &sub.ProblemID, &sub.Code, &sub.Language, &sub.Status, &runtimeMs, &sub.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get submission by id: %w", err)
	}
	if runtimeMs.Valid {
		rt := int(runtimeMs.Int64)
		sub.RuntimeMs = &rt
	}
	return &sub, nil
}

// UpdateSubmissionStatus updates the status and optionally runtime of a submission.
func (s *SQLite) UpdateSubmissionStatus(id int64, status string, runtimeMs *int) error {
	var err error
	if runtimeMs != nil {
		_, err = s.DB.Exec(
			"UPDATE submissions SET status = ?, runtime_ms = ? WHERE id = ?",
			status, *runtimeMs, id,
		)
	} else {
		_, err = s.DB.Exec(
			"UPDATE submissions SET status = ? WHERE id = ?",
			status, id,
		)
	}
	if err != nil {
		return fmt.Errorf("update submission status: %w", err)
	}
	return nil
}

// ListSubmissionsByProblem returns submissions for a problem ordered by creation time.
func (s *SQLite) ListSubmissionsByProblem(problemID int64) ([]Submission, error) {
	rows, err := s.DB.Query(
		"SELECT id, user_id, problem_id, code, language, status, runtime_ms, created_at FROM submissions WHERE problem_id = ? ORDER BY created_at DESC",
		problemID,
	)
	if err != nil {
		return nil, fmt.Errorf("list submissions: %w", err)
	}
	defer rows.Close()

	var subs []Submission
	for rows.Next() {
		var sub Submission
		var runtimeMs sql.NullInt64
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.ProblemID, &sub.Code, &sub.Language, &sub.Status, &runtimeMs, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan submission: %w", err)
		}
		if runtimeMs.Valid {
			rt := int(runtimeMs.Int64)
			sub.RuntimeMs = &rt
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// --- Race Operations ---

// Race represents a competitive coding room.
type Race struct {
	ID         int64             `json:"id"`
	RoomCode   string            `json:"room_code"`
	ProblemID  int64             `json:"problem_id"`
	HostUserID int64             `json:"host_user_id"`
	Status     string            `json:"status"`
	StartTime  *time.Time        `json:"start_time,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	Players    []RaceParticipant `json:"players,omitempty"`
}

// RaceParticipant represents a player in a race.
type RaceParticipant struct {
	RaceID     int64   `json:"race_id"`
	UserID     int64   `json:"user_id"`
	Nickname   string  `json:"nickname,omitempty"`
	Role       string  `json:"role"`
	Status     string  `json:"status"`
	FinishTime *int    `json:"finish_time,omitempty"`
	Verdict    *string `json:"verdict,omitempty"`
}

// CreateRace creates a new race room.
func (s *SQLite) CreateRace(roomCode string, problemID, hostUserID int64) (*Race, error) {
	result, err := s.DB.Exec(
		"INSERT INTO races (room_code, problem_id, host_user_id) VALUES (?, ?, ?)",
		roomCode, problemID, hostUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("create race: %w", err)
	}

	id, _ := result.LastInsertId()

	// Auto-join host as player
	_, err = s.DB.Exec("INSERT INTO race_participants (race_id, user_id, role) VALUES (?, ?, 'player')", id, hostUserID)
	if err != nil {
		return nil, fmt.Errorf("add host to race: %w", err)
	}

	return s.GetRaceByCode(roomCode)
}

// GetRaceByCode returns a race by its room code.
func (s *SQLite) GetRaceByCode(code string) (*Race, error) {
	var race Race
	var startTime sql.NullTime

	err := s.DB.QueryRow(
		"SELECT id, room_code, problem_id, host_user_id, status, start_time, created_at FROM races WHERE room_code = ?",
		code,
	).Scan(&race.ID, &race.RoomCode, &race.ProblemID, &race.HostUserID, &race.Status, &startTime, &race.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get race: %w", err)
	}

	if startTime.Valid {
		race.StartTime = &startTime.Time
	}

	// Get participants
	race.Players, _ = s.GetRaceParticipants(race.ID)

	return &race, nil
}

// GetRaceParticipants returns all participants in a race.
func (s *SQLite) GetRaceParticipants(raceID int64) ([]RaceParticipant, error) {
	rows, err := s.DB.Query(`
		SELECT rp.race_id, rp.user_id, u.nickname, rp.role, rp.status, rp.finish_time, rp.verdict
		FROM race_participants rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.race_id = ?
	`, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []RaceParticipant
	for rows.Next() {
		var p RaceParticipant
		var finishTime sql.NullInt64
		var verdict sql.NullString
		if err := rows.Scan(&p.RaceID, &p.UserID, &p.Nickname, &p.Role, &p.Status, &finishTime, &verdict); err != nil {
			continue
		}
		if finishTime.Valid {
			ft := int(finishTime.Int64)
			p.FinishTime = &ft
		}
		if verdict.Valid {
			p.Verdict = &verdict.String
		}
		players = append(players, p)
	}
	return players, nil
}

// JoinRace adds a user to a race as a player.
func (s *SQLite) JoinRace(raceID, userID int64) error {
	_, err := s.DB.Exec("INSERT OR IGNORE INTO race_participants (race_id, user_id, role) VALUES (?, ?, 'player')", raceID, userID)
	return err
}

// JoinRaceAsSpectator adds a user to a race as a spectator.
func (s *SQLite) JoinRaceAsSpectator(raceID, userID int64) error {
	_, err := s.DB.Exec("INSERT OR IGNORE INTO race_participants (race_id, user_id, role) VALUES (?, ?, 'spectator')", raceID, userID)
	return err
}

// UpdateRaceStatus updates the race status and optionally start time.
func (s *SQLite) UpdateRaceStatus(raceID int64, status string, startTime *time.Time) error {
	if startTime != nil {
		_, err := s.DB.Exec("UPDATE races SET status = ?, start_time = ? WHERE id = ?", status, startTime, raceID)
		return err
	}
	_, err := s.DB.Exec("UPDATE races SET status = ? WHERE id = ?", status, raceID)
	return err
}

// UpdateRaceParticipant updates a participant's status in a race.
func (s *SQLite) UpdateRaceParticipant(raceID, userID int64, status string, finishTime *int, verdict *string) error {
	_, err := s.DB.Exec(
		"UPDATE race_participants SET status = ?, finish_time = ?, verdict = ? WHERE race_id = ? AND user_id = ?",
		status, finishTime, verdict, raceID, userID,
	)
	return err
}
