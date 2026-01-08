-- Migration: Add ELO Rating System
-- Date: 2026-01-09
-- Description: Adds ELO rating tracking for users and race history

-- Add ELO rating column to users table
ALTER TABLE users ADD COLUMN elo_rating INTEGER DEFAULT 1200;

-- Create ELO history table to track rating changes
CREATE TABLE IF NOT EXISTS elo_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    race_id INTEGER NOT NULL,
    old_rating INTEGER NOT NULL,
    new_rating INTEGER NOT NULL,
    rating_change INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (race_id) REFERENCES races(id)
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_elo_history_user ON elo_history(user_id);
CREATE INDEX IF NOT EXISTS idx_elo_history_race ON elo_history(race_id);
CREATE INDEX IF NOT EXISTS idx_users_elo_rating ON users(elo_rating DESC);

-- Update existing users to have default ELO rating (if any exist)
UPDATE users SET elo_rating = 1200 WHERE elo_rating IS NULL;
