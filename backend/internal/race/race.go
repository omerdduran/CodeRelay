package race

import (
	"log"
	"sort"

	"coderelay/backend/internal/elo"
	"coderelay/backend/internal/storage"
	"coderelay/backend/internal/ws"
)

// Service handles race-related business logic
type Service struct {
	store *storage.SQLite
	hub   *ws.Hub
}

// NewService creates a new race service
func NewService(store *storage.SQLite, hub *ws.Hub) *Service {
	return &Service{
		store: store,
		hub:   hub,
	}
}

// CheckRaceCompletion checks if a race is complete and calculates ELO ratings
func (s *Service) CheckRaceCompletion(raceID int64) error {
	// Get all participants
	participants, err := s.store.GetRaceParticipants(raceID)
	if err != nil {
		return err
	}

	// Check if all players have finished (only count players, not spectators)
	allFinished := true
	playerCount := 0
	for _, p := range participants {
		if p.Role == "player" {
			playerCount++
			if p.Status == "racing" {
				allFinished = false
			}
		}
	}

	// If not all players finished, don't calculate ELO yet
	if !allFinished || playerCount == 0 {
		return nil
	}

	log.Printf("race: all %d players finished race %d, calculating ELO", playerCount, raceID)

	// Calculate ELO ratings
	return s.CalculateAndUpdateElo(raceID, participants)
}

// CalculateAndUpdateElo calculates new ELO ratings for race participants
func (s *Service) CalculateAndUpdateElo(raceID int64, participants []storage.RaceParticipant) error {
	// Filter only players who finished with AC (Accepted) verdict
	var finishedPlayers []storage.RaceParticipant
	for _, p := range participants {
		if p.Role == "player" && p.Status == "finished" && p.Verdict != nil && *p.Verdict == "AC" {
			finishedPlayers = append(finishedPlayers, p)
		}
	}

	// If no one finished with AC, no ELO changes
	if len(finishedPlayers) == 0 {
		log.Printf("race: no players finished with AC in race %d, no ELO changes", raceID)
		return nil
	}

	// Sort by finish time (fastest first)
	sort.Slice(finishedPlayers, func(i, j int) bool {
		if finishedPlayers[i].FinishTime == nil {
			return false
		}
		if finishedPlayers[j].FinishTime == nil {
			return true
		}
		return *finishedPlayers[i].FinishTime < *finishedPlayers[j].FinishTime
	})

	// Prepare ELO players
	var eloPlayers []elo.Player
	for rank, p := range finishedPlayers {
		user, err := s.store.GetUserByID(p.UserID)
		if err != nil {
			continue
		}

		eloPlayers = append(eloPlayers, elo.Player{
			UserID: p.UserID,
			Rating: user.EloRating,
			Rank:   rank + 1, // 1-indexed
		})
	}

	// Calculate new ratings
	eloPlayers = elo.CalculateRatings(eloPlayers)

	// Update database and record history
	var eloChanges []map[string]interface{}
	for _, player := range eloPlayers {
		// Update user's ELO rating
		if err := s.store.UpdateUserElo(player.UserID, player.NewRating); err != nil {
			log.Printf("race: failed to update ELO for user %d: %v", player.UserID, err)
			continue
		}

		// Record ELO change in history
		if err := s.store.RecordEloChange(player.UserID, raceID, player.Rating, player.NewRating, player.Rank); err != nil {
			log.Printf("race: failed to record ELO change for user %d: %v", player.UserID, err)
		}

		ratingChange := player.NewRating - player.Rating
		log.Printf("race: user %d: %d -> %d (%+d) [rank %d]",
			player.UserID, player.Rating, player.NewRating, ratingChange, player.Rank)

		eloChanges = append(eloChanges, map[string]interface{}{
			"user_id":       player.UserID,
			"old_rating":    player.Rating,
			"new_rating":    player.NewRating,
			"rating_change": ratingChange,
			"rank":          player.Rank,
		})
	}

	// Broadcast ELO changes to the race room
	race, err := s.getRaceByID(raceID)
	if err == nil && race != nil && s.hub != nil {
		s.hub.BroadcastToRoom(race.RoomCode, ws.TypeRaceEvent, map[string]interface{}{
			"event":       "elo_updated",
			"race_id":     raceID,
			"elo_changes": eloChanges,
		})
	}

	return nil
}

// getRaceByID is a helper to get race by ID
func (s *Service) getRaceByID(raceID int64) (*storage.Race, error) {
	// Query race by ID
	var race storage.Race
	err := s.store.DB.QueryRow(`
		SELECT id, room_code, problem_id, host_user_id, status, start_time, created_at 
		FROM races 
		WHERE id = ?
	`, raceID).Scan(&race.ID, &race.RoomCode, &race.ProblemID, &race.HostUserID, &race.Status, &race.StartTime, &race.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &race, nil
}
