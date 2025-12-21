# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [2024-12-22] - Phase 4: Real-time Features

### Added

#### Backend WebSocket
- `ws/hub.go` - Connection management and broadcasting
- `ws/client.go` - Client with read/write pumps, ping/pong
- WebSocket endpoint: `GET /ws`
- Leaderboard endpoint: `GET /api/leaderboard`

#### Real-time Updates
- Worker broadcasts submission status changes via WebSocket
- Clients receive instant verdict updates

#### Frontend Components
- `useWebSocket.js` - Auto-reconnect WebSocket hook
- `Leaderboard.js` - Live rankings with real-time updates
- `Timer.js` - Countdown timer with warning states

### Dependencies Added
- `github.com/gorilla/websocket` - WebSocket support

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

---

## [2024-12-21] - Phase 2: Basic Frontend

### Added

#### Components
- `NicknameScreen` - User authentication with localStorage
- `CodeEditor` - Monaco Editor with Python syntax
- `ProblemDescription` - Markdown rendering
- `ResultsPanel` - Color-coded verdicts
- `SubmitButton` - Loading state

#### Pages
- Main page (`/`) - Problem list
- Problem page (`/problem/[id]`) - Split layout

### Dependencies Added
- `@monaco-editor/react`, `react-markdown`, `remark-gfm`

---

## [2024-12-21] - Phase 1: Database & Core API

### Added

#### Database Layer
- SQLite with 4 tables
- Embedded migrations
- Seed data (Two Sum + test cases)

#### API Endpoints
- `/healthz`, `/api/problems`, `/api/submissions`
