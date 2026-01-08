package worker

import (
	"context"
	"database/sql"
	"log"
	"time"

	"coderelay/backend/internal/race"
	"coderelay/backend/internal/runner"
	"coderelay/backend/internal/storage"
	"coderelay/backend/internal/ws"
)

// Worker processes queued submissions in the background
type Worker struct {
	store       *storage.SQLite
	runner      *runner.Runner
	hub         *ws.Hub
	raceService *race.Service
	interval    time.Duration
}

// New creates a new submission worker
func New(store *storage.SQLite, hub *ws.Hub) *Worker {
	return &Worker{
		store:       store,
		runner:      runner.New(),
		hub:         hub,
		raceService: race.NewService(store, hub),
		interval:    1 * time.Second,
	}
}

// Start begins processing submissions in the background
func (w *Worker) Start(ctx context.Context) {
	log.Println("worker: starting submission processor")

	// Check for Docker image
	if err := runner.CheckDockerImage(runner.DockerImage); err != nil {
		log.Printf("worker: WARNING - %v", err)
		log.Println("worker: submissions will fail until image is built")
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker: shutting down")
			return
		case <-ticker.C:
			w.processNext()
		}
	}
}

// processNext finds and processes the next queued submission
func (w *Worker) processNext() {
	// Find next queued submission
	sub, err := w.findQueuedSubmission()
	if err != nil || sub == nil {
		return
	}

	log.Printf("worker: processing submission #%d", sub.ID)

	// Update status to running and broadcast
	if err := w.store.UpdateSubmissionStatus(sub.ID, "running", nil); err != nil {
		log.Printf("worker: failed to update status: %v", err)
		return
	}
	w.broadcastUpdate(sub.ID, sub.UserID, sub.ProblemID, "running", nil)

	// Get problem and test cases
	testCases, err := w.store.GetTestCasesByProblemID(sub.ProblemID)
	if err != nil {
		log.Printf("worker: failed to get test cases: %v", err)
		w.store.UpdateSubmissionStatus(sub.ID, "RE", nil)
		w.broadcastUpdate(sub.ID, sub.UserID, sub.ProblemID, "RE", nil)
		return
	}

	problem, err := w.store.GetProblemByID(sub.ProblemID)
	if err != nil {
		log.Printf("worker: failed to get problem: %v", err)
		w.store.UpdateSubmissionStatus(sub.ID, "RE", nil)
		w.broadcastUpdate(sub.ID, sub.UserID, sub.ProblemID, "RE", nil)
		return
	}

	// Configure runner with problem limits
	w.runner.Timeout = time.Duration(problem.TimeLimitMs) * time.Millisecond

	// Run against all test cases
	var totalRuntime time.Duration
	var finalVerdict runner.Verdict = runner.VerdictAC

	for i, tc := range testCases {
		result, err := w.runner.Run(sub.Code, tc.Input)
		if err != nil {
			log.Printf("worker: execution error on test %d: %v", i+1, err)
			finalVerdict = runner.VerdictRE
			break
		}

		totalRuntime += result.Runtime
		verdict := runner.Judge(result, tc.ExpectedOutput)

		log.Printf("worker: test %d/%d: %s (%.0fms)", i+1, len(testCases), verdict, result.Runtime.Seconds()*1000)

		if verdict != runner.VerdictAC {
			finalVerdict = verdict
			break
		}
	}

	// Update final status and broadcast
	runtimeMs := int(totalRuntime.Milliseconds())
	if err := w.store.UpdateSubmissionStatus(sub.ID, string(finalVerdict), &runtimeMs); err != nil {
		log.Printf("worker: failed to update final status: %v", err)
		return
	}
	w.broadcastUpdate(sub.ID, sub.UserID, sub.ProblemID, string(finalVerdict), &runtimeMs)

	log.Printf("worker: submission #%d complete: %s (%dms)", sub.ID, finalVerdict, runtimeMs)

	// Check if this submission is part of a race and update race participant
	w.handleRaceSubmission(sub.ID, sub.UserID, sub.ProblemID, string(finalVerdict), runtimeMs)
}

// broadcastUpdate sends a submission update via WebSocket
func (w *Worker) broadcastUpdate(subID, userID, problemID int64, status string, runtimeMs *int) {
	if w.hub != nil {
		w.hub.BroadcastSubmissionUpdate(subID, userID, problemID, status, runtimeMs)
	}
}

// findQueuedSubmission gets the oldest queued submission
func (w *Worker) findQueuedSubmission() (*storage.Submission, error) {
	row := w.store.DB.QueryRow(`
		SELECT id, user_id, problem_id, code, language, status, runtime_ms, created_at 
		FROM submissions 
		WHERE status = 'queued' 
		ORDER BY created_at ASC 
		LIMIT 1
	`)

	var sub storage.Submission
	var runtimeMs *int
	err := row.Scan(&sub.ID, &sub.UserID, &sub.ProblemID, &sub.Code, &sub.Language, &sub.Status, &runtimeMs, &sub.CreatedAt)
	if err != nil {
		return nil, nil
	}

	return &sub, nil
}

// handleRaceSubmission checks if submission is part of a race and updates race status
func (w *Worker) handleRaceSubmission(subID, userID, problemID int64, verdict string, runtimeMs int) {
	// Find if user is in an active race for this problem
	row := w.store.DB.QueryRow(`
		SELECT r.id, r.room_code, r.start_time 
		FROM races r
		JOIN race_participants rp ON r.id = rp.race_id
		WHERE r.problem_id = ? 
		AND r.status = 'racing'
		AND rp.user_id = ?
		AND rp.role = 'player'
		AND rp.status = 'racing'
	`, problemID, userID)

	var raceID int64
	var roomCode string
	var startTime sql.NullTime
	if err := row.Scan(&raceID, &roomCode, &startTime); err != nil {
		// Not in a race or already finished
		return
	}

	// Calculate finish time from race start
	var finishTime *int
	if startTime.Valid {
		elapsed := int(time.Since(startTime.Time).Milliseconds())
		finishTime = &elapsed
	}

	// Update participant status
	verdictStr := verdict
	if err := w.store.UpdateRaceParticipant(raceID, userID, "finished", finishTime, &verdictStr); err != nil {
		log.Printf("worker: failed to update race participant: %v", err)
		return
	}

	log.Printf("worker: user %d finished race %d with verdict %s in %dms", userID, raceID, verdict, finishTime)

	// Broadcast to race room
	if w.hub != nil {
		w.hub.BroadcastToRoom(roomCode, ws.TypeRaceEvent, map[string]interface{}{
			"event":       "player_finished",
			"user_id":     userID,
			"verdict":     verdict,
			"finish_time": finishTime,
		})
	}

	// Check if race is complete and calculate ELO
	if err := w.raceService.CheckRaceCompletion(raceID); err != nil {
		log.Printf("worker: failed to check race completion: %v", err)
	}
}
