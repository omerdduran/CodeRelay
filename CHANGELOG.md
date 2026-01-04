# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [2024-12-22] - Bug Fixes & Improvements

### Fixed
- Spectator view now correctly displays player code updates. Fixed payload parsing issue where WebSocket messages had string payloads that needed to be parsed as JSON.

### Changed
- `make backend-run` now automatically builds the Docker runner image (`coderelay-runner`) if it doesn't exist, eliminating the need for manual build steps.

---

## [2024-12-22] - Live Race Feature

### Added

#### Database
- `races` table with room_code, status, start_time
- `race_participants` table for player tracking

#### Backend API
- `POST /api/races` - Create race room
- `GET /api/races/{code}` - Get room info
- `POST /api/races/{code}/join` - Join room
- `POST /api/races/{code}/start` - Start race

#### WebSocket Events
- `race_event.player_joined`
- `race_event.countdown`
- `race_event.race_started`

#### Frontend
- `/race` - Create/Join lobby
- `/race/[code]` - Waiting room + Live race view
- Countdown animation (3-2-1)
- Real-time player status

---

## [2024-12-21] - Phase 4: Real-time Features

### Added
- WebSocket hub with client management
- Leaderboard endpoint and live component
- Timer component with warnings

---

## [2024-12-21] - Phase 3: Code Runner

### Added
- Python Docker sandbox
- `runner.go` - Code execution
- `verdict.go` - AC/WA/TLE comparison
- Background worker for submissions

---

## [2024-12-21] - Phase 2: Basic Frontend

### Added
- NicknameScreen, CodeEditor (Monaco), ProblemDescription
- ResultsPanel, SubmitButton
- Main page and problem page

---

## [2024-12-21] - Phase 1: Database & Core API

### Added
- SQLite with users, problems, test_cases, submissions
- REST API endpoints
- Seed data (Two Sum problem)
