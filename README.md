**CodeRelay** is a fast, social way to practice algorithms and data structures. Players join a live room, pick a challenge, write code in the`` browser, and watch rankings update in real time. It’s built for bootcamps, university clubs, engineering teams, and streamers who want competitive energy without setup headaches.

## What makes it fun

- **Live leaderboard:** Rankings update instantly as submissions pass tests (no page refreshes).
    
- **1v1 duels:** Head-to-head matches—first Accepted wins (or best score at the buzzer).
    
- **Time-attack mode:** A visible countdown adds pressure and drama.
    
- **Spectator links:** Share a room URL so viewers can watch the match live.
    
- **Fair play by design:** Code runs server-side in isolated, short-lived containers with strict limits.
    

## How it works

1. Join a room (solo, group, or 1v1) — or share a link for spectators.
    
2. Read the problem, code in the browser, and submit.
    
3. The platform spins up a fresh container, evaluates your solution, and pushes results to everyone in the room in real time.
    
4. When the timer ends, winners and detailed results are announced.
    

## Technology overview

- **Frontend:** React with **Monaco Editor** for a familiar, VS Code–like experience; real-time updates via **WebSockets** (two-way communication—no polling). ([Microsoft GitHub](https://microsoft.github.io/monaco-editor/?utm_source=chatgpt.com "Monaco Editor"))
    
- **Backend:** **Go** for low-latency rooms and judging.
    
- **Database:** **SQLite** (embedded, zero-config) for users, problems, submissions, and scores—great for a lean MVP. ([SQLite](https://sqlite.org/?utm_source=chatgpt.com "SQLite Home Page"))
    
- **Sandboxed execution:** Each submission runs in an **ephemeral Docker** container with CPU/RAM/time limits and **seccomp** syscall filtering; hardening options like **gVisor** can be added in production. ([Docker Documentation](https://docs.docker.com/engine/containers/resource_constraints/?utm_source=chatgpt.com "Resource constraints"))

---

# Developer setup

- **Go:** `go1.22.5` (pinned via `toolchain` in `backend/go.mod`).
- **Node.js:** `22.14.0` (see `.nvmrc` and `frontend/package.json`).
- Install dependencies once with `make frontend-install`.
- Run the API locally with `make backend-run`; stop with `Ctrl+C`.
- Execute backend checks via `make backend-test`, `make backend-fmt`, and `make backend-lint`.
- Build or lint the UI with `make frontend-build` and `make frontend-lint`.
- Bring the full stack up using `make compose-up` (wraps `docker compose up --build`).
- Inspect container readiness with `docker compose ps` and tail service logs via `docker compose logs -f api`.

---

# MVP

## Scope

- Single room, single problem
    
- One language: **Python** (start simple; add more later)
    
- Nickname-only sign-in (cookie/local storage)
    
- **WebSockets** for live leaderboard & submission status (replace polling)
    
- Verdicts: AC, WA, TLE
    
- First AC wins (tie breaks on earliest AC time)
    

## Tech Stack

- **Backend:** Go (REST + WebSocket), single judging worker
    
- **Database:** SQLite
    
- **Frontend:** React + Monaco
    
- **Containers:** Docker Compose (api, frontend); SQLite lives on a volume
    

## Database (initial)

- Migrations: `users`, `problem`, `test_case`, `submission`
    
- Seed: 1 demo problem + 5–10 hidden tests
    
- Basic indexes for smooth reads (e.g., `submission(created_at)`)
    

## Backend (behavior)

- Submission lifecycle: `queued → running → AC/WA/TLE`
    
- Single worker: fetch next queued submission, run tests, persist verdict/runtime
    
- **Live updates:** push submission results & leaderboard over WebSocket
    

## Runner (idea level)

- Auto-launch a fresh container per submission
    
- Feed test input via stdin, capture stdout/stderr
    
- Enforce **time** and **memory** limits; **no network**; cleanup always
    
- Map outcomes to AC/WA/TLE
    
    - Use Docker resource controls and the default **seccomp** allowlist for syscall reduction
## Frontend

- Nickname screen
    
- Problem view (markdown statement + samples)
    
- Editor + Submit button (Monaco)
    
- Live results panel (WebSocket)
    
- Leaderboard + countdown timer (synced to server clock)
## DevOps

- Dockerfiles for api/frontend
    
- Compose with volumes for code, SQLite db, and problem assets
    
- Health checks; one-command local run
    
## Testing

- Happy path: AC submission appears on leaderboard
    
- WA and TLE cases
    
- Concurrent submissions from two browsers (leaderboard updates in real time)
    

---
## Extras (nice to have)

- **Elo rating** for 1v1 and arena modes (simple K-factor to start).
    
- “Code preview” opt-in for players who want to share live progress with spectators.
    
- Weekly mini-tournaments and season stats.

## Team Members
- Bruna Pierobon
- Julia Correia Bindi
- Ömer Duran
- Deniz Can Çalkın
