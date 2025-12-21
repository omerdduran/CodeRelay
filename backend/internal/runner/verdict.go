package runner

import (
	"strings"
)

// Verdict represents the result of judging a submission
type Verdict string

const (
	VerdictAC  Verdict = "AC"  // Accepted
	VerdictWA  Verdict = "WA"  // Wrong Answer
	VerdictTLE Verdict = "TLE" // Time Limit Exceeded
	VerdictRE  Verdict = "RE"  // Runtime Error
)

// Judge compares actual output with expected output
func Judge(result *Result, expectedOutput string) Verdict {
	// Check for timeout
	if result.TimedOut {
		return VerdictTLE
	}

	// Check for runtime error
	if result.Error != nil || result.ExitCode != 0 {
		return VerdictRE
	}

	// Normalize outputs for comparison
	actual := normalizeOutput(result.Output)
	expected := normalizeOutput(expectedOutput)

	// Compare outputs
	if actual == expected {
		return VerdictAC
	}

	return VerdictWA
}

// normalizeOutput cleans up output for fair comparison
func normalizeOutput(s string) string {
	// Trim whitespace from each line and join
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return strings.Join(cleaned, "\n")
}
