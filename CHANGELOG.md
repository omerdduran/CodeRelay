# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [2024-12-21] - Phase 3: Code Runner

### Added

#### Docker Sandbox
- Python 3.11 runner Dockerfile with non-root user
- Resource limits: CPU, memory, process count, timeout
- Security: No network access, read-only filesystem

#### Runner Package
- `runner.go` - Docker-based code execution
- `verdict.go` - Output comparison (AC/WA/TLE/RE)

#### Background Worker
- `worker.go` - Processes queued submissions
- Polls every 1 second for new submissions
- Runs code against all test cases
- Updates submission status with verdict

### Changed
- `main.go` - Now starts worker alongside HTTP server

---

## [2024-12-21] - Phase 2: Basic Frontend

### Added

#### Components
- `NicknameScreen` - User authentication with localStorage persistence
- `CodeEditor` - Monaco Editor with Python syntax highlighting
- `ProblemDescription` - Markdown rendering with react-markdown
- `ResultsPanel` - Submission status display with color-coded verdicts
- `SubmitButton` - Loading state and submission trigger

#### Pages
- Main page (`/`) - Problem list with user info
- Problem page (`/problem/[id]`) - Split layout with editor and results

#### Hooks & Utilities
- `useUser` - User state management with localStorage
- `lib/api.js` - API client for backend communication

#### Styling
- Modern dark theme with glassmorphism effects
- Responsive split layout for problem solving
- Custom scrollbar styling

### Dependencies Added
- `@monaco-editor/react` - Code editor
- `react-markdown` - Markdown rendering
- `remark-gfm` - GitHub Flavored Markdown support

---

## [2024-12-21] - Phase 1: Database & Core API

### Added

#### Database Layer
- SQLite database with `go-sqlite3` driver
- Database schema with 4 tables: `users`, `problems`, `test_cases`, `submissions`
- Embedded migrations using Go's `//go:embed` directive
- Seed data with demo problem "Two Sum" and 5 test cases
- Full CRUD operations for all entities

#### API Endpoints
- `GET /healthz` - Health check endpoint
- `GET /api/problems` - List all problems
- `GET /api/problems/{id}` - Get problem details with sample test cases
- `POST /api/submissions` - Create new code submission
- `GET /api/submissions/{id}` - Get submission status

#### Infrastructure
- CORS middleware for frontend development
- Graceful server shutdown
- Automatic database migrations on startup

#### Testing
- Storage layer unit tests (5 tests)
- API handler unit tests (6 tests)
