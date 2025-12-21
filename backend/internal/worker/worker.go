package worker

import (
	"context"
	"log"
	"time"

	"coderelay/backend/internal/runner"
	"coderelay/backend/internal/storage"
)

// Worker processes queued submissions in the background
type Worker struct {
	store    *storage.SQLite
	runner   *runner.Runner
	interval time.Duration
}

// New creates a new submission worker
func New(store *storage.SQLite) *Worker {
	return &Worker{
		store:    store,
		runner:   runner.New(),
		interval: 1 * time.Second,
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

	// Update status to running
	if err := w.store.UpdateSubmissionStatus(sub.ID, "running", nil); err != nil {
		log.Printf("worker: failed to update status: %v", err)
		return
	}

	// Get problem and test cases
	testCases, err := w.store.GetTestCasesByProblemID(sub.ProblemID)
	if err != nil {
		log.Printf("worker: failed to get test cases: %v", err)
		w.store.UpdateSubmissionStatus(sub.ID, "RE", nil)
		return
	}

	// Get problem for time limit
	problem, err := w.store.GetProblemByID(sub.ProblemID)
	if err != nil {
		log.Printf("worker: failed to get problem: %v", err)
		w.store.UpdateSubmissionStatus(sub.ID, "RE", nil)
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

	// Update final status
	runtimeMs := int(totalRuntime.Milliseconds())
	if err := w.store.UpdateSubmissionStatus(sub.ID, string(finalVerdict), &runtimeMs); err != nil {
		log.Printf("worker: failed to update final status: %v", err)
		return
	}

	log.Printf("worker: submission #%d complete: %s (%dms)", sub.ID, finalVerdict, runtimeMs)
}

// findQueuedSubmission gets the oldest queued submission
func (w *Worker) findQueuedSubmission() (*storage.Submission, error) {
	// Query for oldest queued submission
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
		return nil, nil // No queued submissions
	}

	return &sub, nil
}
