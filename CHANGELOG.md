# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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

### Files Added
- `backend/internal/storage/schema.sql`
- `backend/internal/storage/seed.sql`
- `backend/internal/storage/sqlite.go` (rewritten)
- `backend/internal/storage/sqlite_test.go`
- `backend/internal/server/handlers.go`
- `backend/internal/server/handlers_test.go`
- `backend/internal/server/server.go` (rewritten)

---

<!-- 
Example format:

## [2024-01-15] - Feature Name

### Added
- New features added

### Changed
- Changes to existing functionality

### Fixed
- Bug fixes

### Removed
- Removed features

### Security
- Security updates
-->
