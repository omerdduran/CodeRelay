# ELO Rating System

## Overview

CodeRelay implements a competitive ELO rating system to rank players based on their performance in coding races. The system is inspired by chess ELO ratings but adapted for multi-player competitive programming contests.

## How It Works

### Starting Rating
- All new users start with an ELO rating of **1200**
- This is the default baseline for all players

### Rating Calculation

#### When Ratings Are Updated
- Ratings are calculated **after each race** when all players have finished
- Only players who finish with an **AC (Accepted)** verdict are included in rating calculations
- Spectators do not participate in rating calculations
- Single-player races do not affect ELO (competition is required)

#### Formula
We use the standard ELO formula adapted for multi-player contests:

1. **Expected Score**: Calculated based on rating differences between all players
   ```
   E(A) = 1 / (1 + 10^((RatingB - RatingA) / 400))
   ```

2. **Actual Score**: Based on finishing position
   - 1st place: 1.0
   - Last place: 0.0
   - Middle positions: Linear interpolation

3. **Rating Change**:
   ```
   NewRating = OldRating + K × (ActualScore - ExpectedScore)
   ```
   
   Where K = 32 (K-factor determines volatility)

### Multi-Player Races

For races with 3+ players:
- Each player's expected score is the average of their expected win probability against all opponents
- Actual score is normalized based on rank (1st = 1.0, 2nd = 0.66, 3rd = 0.33, etc.)
- This ensures fair rating distribution across all participants

### Key Features

#### Upsets Give More Points
- If a lower-rated player beats a higher-rated player, they gain more points
- The higher-rated player loses more points for losing to a lower-rated player
- Example: 1000-rated player beating a 1400-rated player gains ~28 points

#### Expected Outcomes Give Fewer Points
- If a higher-rated player beats a lower-rated player (expected outcome), they gain fewer points
- Example: 1400-rated player beating a 1000-rated player gains ~4 points

#### Rating Floor
- Minimum rating is **100**
- Players cannot drop below this threshold

#### Zero-Sum System
- In 2-player races, total rating change is balanced (winner gains what loser loses)
- In multi-player races, ratings are distributed fairly based on performance

## API Endpoints

### Get Global ELO Leaderboard
```http
GET /api/leaderboard?type=elo&limit=100
```

Returns top players sorted by ELO rating.

**Response:**
```json
[
  {
    "rank": 1,
    "user_id": 123,
    "nickname": "CodeMaster",
    "elo_rating": 1456,
    "created_at": "2026-01-09T10:00:00Z"
  }
]
```

### Get User ELO History
```http
GET /api/users/{user_id}/elo-history?limit=50
```

Returns a user's rating change history.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 123,
    "race_id": 45,
    "old_rating": 1200,
    "new_rating": 1228,
    "rating_change": 28,
    "rank": 1,
    "created_at": "2026-01-09T10:00:00Z"
  }
]
```

## Database Schema

### Users Table
```sql
ALTER TABLE users ADD COLUMN elo_rating INTEGER DEFAULT 1200;
```

### ELO History Table
```sql
CREATE TABLE elo_history (
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
```

## WebSocket Events

### ELO Updated Event
When a race completes and ratings are calculated, all clients in the race room receive:

```json
{
  "type": "race_event",
  "payload": {
    "event": "elo_updated",
    "race_id": 45,
    "elo_changes": [
      {
        "user_id": 123,
        "old_rating": 1200,
        "new_rating": 1228,
        "rating_change": 28,
        "rank": 1
      }
    ]
  }
}
```

## Frontend Integration

### Leaderboard Component
The `Leaderboard` component supports two modes:

1. **ELO Mode** (default): Shows global rankings by ELO rating
2. **Problem Mode**: Shows problem-specific rankings by solve time

Users can toggle between modes using the UI controls.

### Props
```jsx
<Leaderboard 
  problemId={null}  // null for global, or specific problem ID
  showElo={true}    // enable ELO mode toggle
/>
```

## Implementation Details

### Backend
- **Package**: `backend/internal/elo`
- **Service**: `backend/internal/race` (handles race completion)
- **Storage**: ELO functions in `backend/internal/storage/sqlite.go`

### Race Completion Flow
1. Player submits code
2. Worker processes submission
3. If submission is part of a race, update race participant status
4. Check if all players have finished
5. Calculate ELO ratings for all finishers with AC verdict
6. Update user ratings and record history
7. Broadcast ELO changes via WebSocket

### Testing
Comprehensive unit tests in `backend/internal/elo/elo_test.go` cover:
- 2-player scenarios
- Multi-player scenarios
- Upsets (unexpected outcomes)
- Expected outcomes
- Rating floor enforcement
- Score calculation accuracy

## Future Enhancements

Potential improvements for the ELO system:

1. **Seasonal Resets**: Reset ratings at the start of each season
2. **Rating Decay**: Inactive players gradually lose rating
3. **Provisional Ratings**: New players have provisional ratings until they complete N races
4. **Rating Tiers**: Bronze, Silver, Gold, Platinum, Diamond ranks
5. **K-Factor Adjustments**: Lower K-factor for high-rated players (less volatility)
6. **Problem Difficulty Weighting**: Harder problems affect rating more

## References

- [ELO Rating System - Wikipedia](https://en.wikipedia.org/wiki/Elo_rating_system)
- [Chess ELO Calculation](https://www.chess.com/terms/elo-rating-chess)
- [Multi-Player ELO Systems](https://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details)
