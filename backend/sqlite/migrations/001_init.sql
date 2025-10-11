PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS problems (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    statement TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS matches (
    id TEXT PRIMARY KEY,
    problem_id TEXT NOT NULL REFERENCES problems(id),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'in_progress', 'finished')
    ),
    config JSON,
    started_at DATETIME,
    finished_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS match_participants (
    match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    seat INTEGER NOT NULL CHECK (seat IN (1, 2)),
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (match_id, player_id),
    UNIQUE (match_id, seat)
);

CREATE TABLE IF NOT EXISTS submissions (
    id TEXT PRIMARY KEY,
    match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    language TEXT NOT NULL DEFAULT 'python',
    verdict TEXT NOT NULL DEFAULT 'pending' CHECK (
        verdict IN ('pending', 'AC', 'WA', 'TLE', 'RE')
    ),
    runtime_ms INTEGER,
    wall_time_ms INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS submission_audit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    submission_id TEXT NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    stage TEXT NOT NULL CHECK (
        stage IN ('compile', 'execute', 'post_process')
    ),
    status TEXT NOT NULL CHECK (
        status IN ('pending', 'running', 'success', 'failed')
    ),
    detail TEXT,
    stdout TEXT,
    stderr TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_matches_problem ON matches(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_match ON submissions(match_id);
CREATE INDEX IF NOT EXISTS idx_submissions_player ON submissions(player_id);
CREATE INDEX IF NOT EXISTS idx_submission_audit_submission ON submission_audit(submission_id);
