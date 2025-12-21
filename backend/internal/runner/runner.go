package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	// Docker image for Python execution
	DockerImage = "coderelay-runner"

	// Default limits
	DefaultTimeout    = 5 * time.Second
	DefaultMemoryMB   = 256
	DefaultCPUPercent = 50
)

// Result contains the execution result
type Result struct {
	Output   string
	Runtime  time.Duration
	ExitCode int
	TimedOut bool
	Error    error
}

// Runner executes code in a Docker container
type Runner struct {
	Image    string
	Timeout  time.Duration
	MemoryMB int
	CPUPct   int
}

// New creates a new Runner with default settings
func New() *Runner {
	return &Runner{
		Image:    DockerImage,
		Timeout:  DefaultTimeout,
		MemoryMB: DefaultMemoryMB,
		CPUPct:   DefaultCPUPercent,
	}
}

// Run executes Python code with the given input
func (r *Runner) Run(code, input string) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout+2*time.Second)
	defer cancel()

	// Build docker command
	args := []string{
		"run",
		"--rm",              // Remove container after execution
		"--network", "none", // No network access
		"--read-only",                               // Read-only filesystem
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=64m", // Writable /tmp
		fmt.Sprintf("--memory=%dm", r.MemoryMB),
		fmt.Sprintf("--cpus=%.2f", float64(r.CPUPct)/100.0),
		"--pids-limit", "50", // Limit processes
		"-i", // Interactive (for stdin)
		r.Image,
		"python3", "-c", code,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Set up stdin
	cmd.Stdin = strings.NewReader(input)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute
	start := time.Now()
	err := cmd.Run()
	runtime := time.Since(start)

	result := &Result{
		Output:  strings.TrimSpace(stdout.String()),
		Runtime: runtime,
	}

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.Error = fmt.Errorf("execution timed out after %v", r.Timeout)
		return result, nil
	}

	// Check for execution error
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			// Include stderr in output for debugging
			if stderr.Len() > 0 {
				result.Output = strings.TrimSpace(stderr.String())
			}
		} else {
			result.Error = err
		}
	}

	return result, nil
}

// CheckDockerImage verifies the runner image exists
func CheckDockerImage(image string) error {
	cmd := exec.Command("docker", "image", "inspect", image)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker image '%s' not found, run: docker build -t %s backend/runner/", image, image)
	}
	return nil
}
